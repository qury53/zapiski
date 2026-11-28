package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	userv1 "shared/pkg/proto/user/v1"
	"time"

	grpcProm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qury53/GC_learn/platform/auth"
	"github.com/qury53/GC_learn/platform/interceptor"
	"github.com/qury53/GC_learn/platform/logger"
	prometheus_metrics "github.com/qury53/GC_learn/platform/metrics"
	pgx_pool "github.com/qury53/GC_learn/platform/postgres"
	"github.com/qury53/GC_learn/platform/redis"
	user_repository "github.com/qury53/GC_learn/user/internal/repository"
	user_service "github.com/qury53/GC_learn/user/internal/service"
	user_transport "github.com/qury53/GC_learn/user/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "user"
)

type App struct {
	log          *logger.Logger
	grpcServer   *grpc.Server
	metricsSrv   *http.Server
	metricsReg   *prometheus.Registry
	port         string
	shutdownTime time.Duration
	closers      []func()
}

func New(
	ctx context.Context,
	config *Config,
	log *logger.Logger,
) (*App, error) {
	app := &App{
		log:          log,
		port:         config.GRPCPort,
		shutdownTime: config.ShutdownTime,
	}

	redisConfig := redis.NewConfigMust()
	postgresConfig := pgx_pool.NewConfigMust()

	//инициализация клиентов, пуллов и тд
	redisClient, err := redis.NewClient(ctx, redisConfig)
	if err != nil {
		log.Error("failed to initialize redis client", slog.Any("error", err))
		app.Close()
		return nil, fmt.Errorf("redis client initialization failed: %w", err)
	}

	app.closers = append(app.closers, func() { _ = redisClient.Close() })

	pool, err := pgx_pool.NewPool(ctx, postgresConfig)

	if err != nil {
		log.Error("failed to inizitalize pool", slog.Any("error", err))
		app.Close()
		return nil, fmt.Errorf("pool initialization failed: %w", err)
	}
	app.closers = append(app.closers, pool.Close)

	// инициализация фич

	srvMetrics := grpcProm.NewServerMetrics(
		grpcProm.WithServerHandlingTimeHistogram(),
	)

	reg := prometheus_metrics.NewRegistry()
	reg.MustRegister(srvMetrics)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
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

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.TimeoutUnaryInterceptor(15*time.Second),
		interceptor.RecoveryInterceptor(log),
		auth.UnaryServerInterceptor(authValidator, authPolicy),
		srvMetrics.UnaryServerInterceptor(),
		interceptor.LoggingUnaryInterceptor(log, serviceName),
	))

	userRepo := user_repository.NewUserRepository(pool)
	userService := user_service.NewUserService(log, userRepo, redisClient)
	userTransport := user_transport.NewUserTransport(log, userService)

	userv1.RegisterUserServer(grpcServer, userTransport)
	srvMetrics.InitializeMetrics(grpcServer)

	if config.Env == "dev" {
		reflection.Register(grpcServer)
	}

	app.grpcServer = grpcServer

	return app, nil
}

func (a *App) Close() {
	// Закрываем ресурсы в обратном порядке (LIFO)
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
		if err := a.metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error(
				"metrics server error",
				slog.Any("error", err),
			)

			errChan <- err
		}
	}()

	go func() {
		a.log.Info("user запущен", slog.String("port", a.port))
		if err := a.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("остановка user")
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
			a.log.Info("user сервис остановлен корректно")
		case <-shutdownCtx.Done():
			a.log.Warn("превышен shutdownTime")
			a.grpcServer.Stop()
		}

	case err := <-errChan:
		a.log.Error("user service error", slog.Any("error", err))
		return err
	}

	return nil
}
