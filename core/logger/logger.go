package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLevel = atomic.NewInt32(int32(LDebug))
)

func SetLevel(level Level) {
	globalLevel.Store(int32(level))
}

type (
	Level = zapcore.Level
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

func enabler() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(globalLevel.Load())
	}
}

// do not migrate this to application OnStart hooks
// as we need a correctly configured logger ASAP
func init() {
	zap.ReplaceGlobals(New())
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
		cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(globalLevel.Load()))
		logger, err = cfg.Build()
	case "PROD":
		logger, err = zap.NewProduction()
	default:
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(globalLevel.Load()))
		logger, err = cfg.Build()
	}
	if err != nil {
		panic(stacktrace.Propagate(err, "unable to initialise logger"))
	}
	return logger
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
