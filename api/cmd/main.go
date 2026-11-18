package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os/signal"
	agentv1 "shared/pkg/proto/agent/v1"
	ssov1 "shared/pkg/proto/sso/v1"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env", err)
	}

	connSSO, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	ssoClient := ssov1.NewAuthServiceClient(connSSO)

	log.Println("sso client created")

	connAgent, err := grpc.NewClient("localhost:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	agentClient := agentv1.NewAgentServiceClient(connAgent)

	log.Println("agent client created")

	go func() {
		log.Printf("starting HTTP server")
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("serve error", slog.String("error", err.Error()))
			cancel()
		}
	}()

	<-ctx.Done()

	log.Printf("shutting down HTTP server")
}
