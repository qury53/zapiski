package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

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
	"google.golang.org/grpc"
)

type App struct {
	log        *logger.Logger
	grpcServer *grpc.Server
	config     *Config
	port       string
	closers    []func()
}

func New(
	ctx context.Context,
	config *Config,
	log *logger.Logger,
) (*App, error) {
	app := App{
		log:  log,
		port: config.GRPCPort,
	}

	serviceConfig := config.NewConfigMust()
	redisConfig := redis.NewConfigMust()
	postgresConfig := pgx_pool.NewConfigMust()
	smtpConfig := email_smtp_repository.NewConfigMust()
	authConfig := auth_service.NewConfigMust()
	smtpConfig := email_smtp_repository.NewConfigMust()
	oauthGoogleConfig := oauth_google_repository.NewConfigMust()
	encryptorConfig := core_encrypt.NewConfigMust()

	//инициализация клиентов, пуллов и тд
	redisClient, err := redis.NewClient(ctx, redisConfig)
	if err != nil {
		log.Error("failed to initialize redis client", slog.Any("error", err))
	}

	oauthGoogleClient := oauth_google_repository.NewClient(oauthGoogleConfig)
	oauthClients := make(map[domain.Provider]core_oauth.Client)
	oauthClients[domain.ProviderGoogle] = oauthGoogleClient

	pool, err := pgx_pool.NewPool(ctx, postgresConfig)
	if err != nil {
		log.Error("failed to inizitalize pool", slog.Any("error", err))
	}

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

	// инициализация gRPC сервера
	list, err := net.Listen("tcp", fmt.Sprintf(":%s", serviceConfig.GRPCPort))
	if err != nil {
		log.Error(
			"порт занят",
			slog.String("port", config.GRPCPort),
			slog.Any("error", err),
		)
	}

	grpcServer := grpc.NewServer()
}
