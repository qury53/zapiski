package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger //встраиваем логгер прям в структуру, чтобы мы могли писать Logger.Info,а не Logger.zaplogger.Info,
	// также это позволяет не терять другие поля структуры при переопределении методов из коробки
	file *os.File
}

type loggerContextKey struct {
}

var (
	key = loggerContextKey{}
)

func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("create folder for log level: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000")
	logFIlePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFIlePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())
	return &Logger{zapLogger, logFile}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close logger", err)
	}
}
