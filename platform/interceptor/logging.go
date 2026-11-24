package interceptor

import (
	"context"
	"log/slog"
	"time"

	"github.com/qury53/GC_learn/platform/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor логгирует унарный rpc запрос
func LoggingUnaryInterceptor(
	log *logger.Logger,
	service string,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		log.Info(
			"started handling request",
			slog.String("service", service),
			slog.String("method", info.FullMethod),
			slog.Time("started handling", start),
		)

		resp, err := handler(ctx, req)

		log.Info(
			"finished handling request",
			slog.String("service", service),
			slog.String("method", info.FullMethod),
			slog.Any("status code", status.Code(err)),
			slog.Duration("time spent", time.Since(start)),
		)

		return resp, err
	}
}
