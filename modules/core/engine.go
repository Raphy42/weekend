package core

import (
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler"
)

type Engine struct {
	bus        message.Bus
	background *scheduler.Pipeline
}
