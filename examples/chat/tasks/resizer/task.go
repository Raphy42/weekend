package resizer

import (
	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/task"
)

var (
	Name = dep.Name("image", "resizer")
	Task = task.Of(Name, ResizeImage)
)
