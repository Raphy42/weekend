//go:build gorm.postgres

package gorm

import (
	"gorm.io/gorm"

	"github.com/Raphy42/weekend/core/app"
)

func postgresDialector() func(dsn string) gorm.Dialector {
	return func(dsn string) gorm.Dialector {
		return postgres.Open(dsn)
	}
}

func init() {
	app.RegisterOnStartHook(func() error {
		globalDriver.Register("postgres", postgresDialector())
	})
	app.RegisterOnStartHook(func() error {
		globalDriver.Register("postgresql", postgresDialector())
	})
}
