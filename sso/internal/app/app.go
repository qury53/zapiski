package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	ssov1 "shared/pkg/proto/sso/v1"
	"time"

	"github.com/qury53/GC_learn/platform/logger"
	pgx_pool "github.com/qury53/GC_learn/platform/postgres"
	"github.com/qury53/GC_learn/platform/redis"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_encrypt "github.com/qury53/GC_learn/sso/internal/core/encrypt"
	core_oauth "github.com/qury53/GC_learn/sso/internal/core/oauth"
	email_smtp_repository "github.com/qury53/GC_learn/sso/internal/repository/email"
	oauth_postgres_repository "github.com/qury53/GC_learn/sso/internal/repository/oauth"
	oauth_google_repository "github.com/qury53/GC_learn/sso/internal/repository/oauth/google"
	sessions_postgres_repository "github.com/qury53/GC_learn/sso/internal/repository/session"
	users_postgres_repository "github.com/qury53/GC_learn/sso/internal/repository/user"
	auth_service "github.com/qury53/GC_learn/sso/internal/service/auth"
	auth_transport "github.com/qury53/GC_learn/sso/internal/transport/auth"
	"google.golang.org/grpc"
)

type App struct {
	log          *logger.Logger
	grpcServer   *grpc.Server
	shutdownTime time.Duration
	port         string
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
	smtpConfig := email_smtp_repository.NewConfigMust()
	authConfig := auth_service.NewConfigMust()
	oauthGoogleConfig := oauth_google_repository.NewConfigMust()
	encryptorConfig := core_encrypt.NewConfigMust()

	//инициализация клиентов, пуллов и тд
	redisClient, err := redis.NewClient(ctx, redisConfig)
	if err != nil {
		log.Error("failed to initialize redis client", slog.Any("error", err))
		app.Close()
		return nil, fmt.Errorf("redis client initialization failed: %w", err)
	}

	app.closers = append(app.closers, func() { _ = redisClient.Close() })

	oauthGoogleClient := oauth_google_repository.NewClient(oauthGoogleConfig)
	oauthClients := make(map[domain.Provider]core_oauth.Client)
	oauthClients[domain.ProviderGoogle] = oauthGoogleClient

	pool, err := pgx_pool.NewPool(ctx, postgresConfig)

	if err != nil {
		log.Error("failed to inizitalize pool", slog.Any("error", err))
		app.Close()
		return nil, fmt.Errorf("pool initialization failed: %w", err)
	}

	app.closers = append(app.closers, pool.Close)

	smtp := email_smtp_repository.NewSmtpRepository(smtpConfig)

	encryptor := core_encrypt.NewEncryptor(encryptorConfig)

	// инициализация фич
	sessionRepository := sessions_postgres_repository.NewSessionRepository(pool)
	userRepositry := users_postgres_repository.NewUserRepository(pool)
	oauthRepository := oauth_postgres_repository.NewOAuthRepository(pool)

	authService := auth_service.NewAuthService(
		log,
		authConfig,
		userRepositry,
		sessionRepository,
		oauthRepository,
		encryptor,
		smtp,
		redisClient,
		oauthClients,
	)

	authTransport := auth_transport.NewAuthTransport(log, authService, authService)
	grpcServer := grpc.NewServer()

	ssov1.RegisterAuthServiceServer(grpcServer, authTransport)
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
		a.log.Info("sso запущен", slog.String("port", a.port))
		if err := a.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("остановка sso")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTime)
		defer cancel()

		closed := make(chan struct{})

		go func() {
			a.grpcServer.GracefulStop()
			close(closed)
		}()
		select {
		case <-closed:
			a.log.Info("sso сервис остановлен корректно")
		case <-shutdownCtx.Done():
			a.log.Warn("превышен shutdownTime")
			a.grpcServer.Stop()
		}

	case err := <-errChan:
		a.log.Error("sso service error", slog.Any("error", err))
		return err
	}

	return nil
}
