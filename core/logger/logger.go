package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	KFile    = "file"
	KLine    = "line"
	LDebug   = zapcore.DebugLevel
	LInfo    = zapcore.InfoLevel
	LWarning = zapcore.WarnLevel
	LError   = zapcore.ErrorLevel
	LFatal   = zapcore.FatalLevel
)

func init() {
	logMode := os.Getenv("WEEKEND_LOG_MODE")
	if logMode == "" {
		logMode = "DEV"
	}
	var logger *zap.Logger
	var err error

	switch logMode {
	case "DEV":
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = cfg.Build()
	case "PROD":
		logger, err = zap.NewProduction()
	default:
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = cfg.Build()
	}
	if err != nil {
		panic(stacktrace.Propagate(err, "unable to initialise root logger"))
	}

	zap.ReplaceGlobals(logger)
}

func New(opts ...Option) *zap.Logger {
	options := newLoggerOptions()
	options.apply(opts...)
	//
	//caller := runtime.CallerName(options.SkipCallFrame)
	//
	//name := caller
	//if options.Name != "" {
	//	name = options.Name
	//}

	return zap.L()
}

func ctxDecorator(ctx context.Context) []Option {
	opts := make([]Option, 0)
	deadline, ok := ctx.Deadline()
	if ok {
		opts = append(opts,
			Decorate(func(s string) string {
				return fmt.Sprintf("%s!", s)
			}),
			Fields(zap.Time("context.deadline", deadline)),
		)
	}

	opts = append(opts, SkipCallFrame(2))
	return opts
}

func FromContext(ctx context.Context, opts ...Option) *zap.Logger {
	//todo context aware metadata retrieval
	// - http
	// - business
	// - domain specific
	opts = append(opts, ctxDecorator(ctx)...)
	return New(opts...)
}

func Flush() {
	_ = zap.L().Sync()
}
