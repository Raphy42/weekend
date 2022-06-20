//go:build app.unsafe

package app

import "context"

func (a *App) UnsafeExecute(ctx context.Context, fn any) error {
	return a.container.UnsafeExecute(ctx, fn)
}
