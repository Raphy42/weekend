package app

import "context"

type Context struct {
	context.Context
	Cancel context.CancelFunc
}
