package main

import (
	"context"
	"os/signal"
	"syscall"
)

const (
	serviceName = "sso"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

}
