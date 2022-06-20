//go:build gorm.sqlite

package gorm

import (
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"

	"github.com/Raphy42/weekend/core"
)

func sqliteDialector() func(dsn string) gorm.Dialector {
	return func(dsn string) gorm.Dialector {
		return sqlite.Open(dsn)
	}
}

func init() {
	core.RegisterOnStartHook(func() error {
		globalDriver.Register("sqlite", postgresDialector())
	})
}
