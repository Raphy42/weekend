package logger

import "go.uber.org/zap"

type LoggerOptions struct {
	SkipCallFrame   int
	Name            string
	NamingDecorator func(string) string
	Fields          []zap.Field
}

type LoggerOption func(options *LoggerOptions)

func (l *LoggerOptions) apply(options ...LoggerOption) {
	for _, opt := range options {
		opt(l)
	}
}

func newLoggerOptions() LoggerOptions {
	return LoggerOptions{
		SkipCallFrame:   1,
		NamingDecorator: defaultNamingDecorator,
		Fields:          make([]zap.Field, 0),
	}
}

func defaultNamingDecorator(s string) string {
	return s
}

func SkipCallFrame(count int) LoggerOption {
	return func(options *LoggerOptions) {
		options.SkipCallFrame = count
	}
}

func Named(name string) LoggerOption {
	return func(options *LoggerOptions) {
		options.Name = name
	}
}

func Decorate(fn func(string) string) LoggerOption {
	return func(options *LoggerOptions) {
		options.NamingDecorator = fn
	}
}

func Fields(fields ...zap.Field) LoggerOption {
	return func(options *LoggerOptions) {
		options.Fields = append(options.Fields, fields...)
	}
}
