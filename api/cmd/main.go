package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	ssov1 "shared/pkg/proto/sso/v1"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	adapters_sso "github.com/qury53/GC_learn/api/internal/adapters/sso"
	service_sso "github.com/qury53/GC_learn/api/internal/service/sso"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
	platform_logger "github.com/qury53/GC_learn/platform/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("failed to load env", err)
	}

	loggerCfg := platform_logger.NewConfigMust()
	appLog := platform_logger.New(loggerCfg, "api")

	ssoAddr := os.Getenv("SSO_GRPC_ADDR")
	if ssoAddr == "" {
		ssoAddr = "localhost:50051"
	}
	connSSO, err := grpc.NewClient(ssoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		appLog.Error("sso grpc connect", slog.Any("error", err))
		os.Exit(1)
	}
	defer connSSO.Close()

	ssoProtoClient := ssov1.NewAuthServiceClient(connSSO)
	ssoAdapter := adapters_sso.NewGRPCClient(ssoProtoClient, appLog)
	ssoService := service_sso.NewAuthService(ssoAdapter, appLog)
	ssoHandler := transport_sso.NewAuthHandler(ssoService, appLog)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", ssoHandler.Register)
	mux.HandleFunc("POST /auth/login", ssoHandler.Login)
	mux.HandleFunc("POST /auth/refresh", ssoHandler.Refresh)
	mux.HandleFunc("POST /auth/logout", ssoHandler.Logout)
	mux.HandleFunc("POST /auth/verify-email", ssoHandler.VerifyEmail)
	mux.HandleFunc("POST /auth/oauth/connect", ssoHandler.GetOAuthConnectUrl)
	mux.HandleFunc("GET /auth/oauth/callback", ssoHandler.HandleOAuthCallback)

	httpAddr := os.Getenv("API_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}
	srv := &http.Server{Addr: httpAddr, Handler: mux}

	go func() {
		appLog.Info("starting HTTP server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLog.Error("serve error", slog.Any("error", err))
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLog.Error("shutdown error", slog.Any("error", err))
	}
	appLog.Info("HTTP server stopped")
}
