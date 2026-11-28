package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	v1 "shared/pkg/proto/agent/v1"
	"time"

	grpcProm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ai_provider_smartapi "github.com/qury53/GC_learn/agent/internal/core/pkg/ai_provider/smartapi"
	chat_postgres_repository "github.com/qury53/GC_learn/agent/internal/repository/chat"
	message_postgres_repository "github.com/qury53/GC_learn/agent/internal/repository/message"
	agent_service "github.com/qury53/GC_learn/agent/internal/service/agent"
	agent_transport "github.com/qury53/GC_learn/agent/internal/transport/agent"
	"github.com/qury53/GC_learn/platform/auth"
	"github.com/qury53/GC_learn/platform/interceptor"
	"github.com/qury53/GC_learn/platform/logger"
	prometheus_metrics "github.com/qury53/GC_learn/platform/metrics"
	pgx_pool "github.com/qury53/GC_learn/platform/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "agent"
)

type App struct {
	grpcServer   *grpc.Server
	log          *logger.Logger
	metricsReg   *prometheus.Registry
	metricsSrv   *http.Server
	port         string
	closers      []func()
	shutdownTime time.Duration
}

func New(
	ctx context.Context,
	log *logger.Logger,
	config *Config,
) (*App, error) {
	app := &App{
		log:          log,
		port:         config.GRPC_port,
		shutdownTime: config.ShutdownTime,
	}

	providerConfig := ai_provider_smartapi.NewConfigMust()
	// deprecated gc provider
	// provider := ai_provider_gc.NewGCProvider(config)
	provider := ai_provider_smartapi.NewSmartAPIProvider(providerConfig, log)

	postgresConfig := pgx_pool.NewConfigMust()
	pool, err := pgx_pool.NewPool(ctx, postgresConfig)
	if err != nil {
		log.Error("failed to initialize postgres pool", slog.Any("error", err))
		app.Close()

		return nil, fmt.Errorf("postgres pool initialization failed: %w", err)
	}

	app.closers = append(app.closers, pool.Close)

	chatRepository := chat_postgres_repository.NewChatRepository(pool)
	messageRepository := message_postgres_repository.NewMessageRepository(pool)

	service := agent_service.NewAgentService(log, provider, chatRepository, messageRepository)
	transport := agent_transport.NewAgentTransport(service, log)

	srvMetrics := grpcProm.NewServerMetrics(
		grpcProm.WithServerHandlingTimeHistogram(),
	)

	reg := prometheus_metrics.NewRegistry()
	reg.MustRegister(srvMetrics)

	app.metricsReg = reg

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("healthy"))
	})

	app.metricsSrv = &http.Server{
		Addr:              net.JoinHostPort("", config.MetricsPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	authConfig := auth.NewConfigMust()
	authValidator := auth.NewValidator(authConfig.Secret, authConfig.Issuer)
	authPolicy := auth.NewPolicy()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.TimeoutUnaryInterceptor(30*time.Second),
			interceptor.RecoveryInterceptor(log),
			auth.UnaryServerInterceptor(authValidator, authPolicy),
			srvMetrics.UnaryServerInterceptor(),
			interceptor.LoggingUnaryInterceptor(log, serviceName),
		),
	)

	v1.RegisterAgentServiceServer(s, transport)
	srvMetrics.InitializeMetrics(s)
	if config.Env == "dev" {
		reflection.Register(s)
	}

	app.grpcServer = s

	return app, nil
}

func (a *App) Close() {
	for i := len(a.closers) - 1; i >= 0; i-- {
		a.closers[i]()
	}

	a.closers = nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.Close()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		a.log.Error(
			"порт занят",
			slog.String("port", a.port),
			slog.Any("error", err),
		)

		return fmt.Errorf("listen failed: %w", err)
	}

	errChan := make(chan error, 1)

	go func() {
		a.log.Info(
			"запуск http сервера метрик",
			slog.String("addr", a.metricsSrv.Addr),
		)

		if err := a.metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error(
				"metrics error",
				slog.Any("error", err),
			)

			errChan <- err
		}
	}()

	go func() {
		a.log.Info("agent запущен", slog.String("port", a.port))
		if err := a.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("остановка agent")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTime)
		defer cancel()

		closed := make(chan struct{})

		go func() {
			_ = a.metricsSrv.Shutdown(shutdownCtx)
			a.grpcServer.GracefulStop()
			close(closed)
		}()

		select {
		case <-closed:
			a.log.Info("agent сервис остановлен корректно")
		case <-shutdownCtx.Done():
			a.log.Warn("превышен shutdownTime")
			a.grpcServer.Stop()
		}

	case err := <-errChan:
		a.log.Error("agent service error", slog.Any("error", err))
		return err
	}

	return nil
}
