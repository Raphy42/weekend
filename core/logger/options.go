package logger

import "go.uber.org/zap"

type Options struct {
	SkipCallFrame   int
	Name            string
	NamingDecorator func(string) string
	Fields          []zap.Field
}

type Option func(options *Options)

func (l *Options) apply(options ...Option) {
	for _, opt := range options {
		opt(l)
	}
}

func newLoggerOptions() Options {
	return Options{
		SkipCallFrame:   1,
		NamingDecorator: defaultNamingDecorator,
		Fields:          make([]zap.Field, 0),
	}
}

func defaultNamingDecorator(s string) string {
	return s
}

func SkipCallFrame(count int) Option {
	return func(options *Options) {
		options.SkipCallFrame = count
	}
}

func Named(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func Decorate(fn func(string) string) Option {
	return func(options *Options) {
		options.NamingDecorator = fn
	}
}

func Fields(fields ...zap.Field) Option {
	return func(options *Options) {
		options.Fields = append(options.Fields, fields...)
	}
}
