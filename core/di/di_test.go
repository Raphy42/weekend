package di

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummySystem struct {
	data int
}

var testModule = Declare("test.module",
	Providers(func() context.Context {
		return context.Background()
	}),
	Invoke(func(d *time.Duration) error {
		return nil
	}),
	Invoke(func(ctx context.Context, system *dummySystem) error {
		return nil
	}),
	Providers(func(ctx context.Context) (*dummySystem, error) {
		return &dummySystem{data: 0x1337}, nil
	}),
)

func Test_DI(t *testing.T) {
	t.Setenv("WEEKEND_LOG_MODE", "PROD")
	a := assert.New(t)

	rootCtx := context.Background()
	testCtx, cancel := context.WithDeadline(rootCtx, time.Now().Add(time.Millisecond*250))
	defer cancel()

	container := NewContainer("test")
	a.Error(container.Use(testModule).start(testCtx))
}
