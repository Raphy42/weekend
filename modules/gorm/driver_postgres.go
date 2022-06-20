//go:build gorm.postgres

package gorm

import (
	"gorm.io/gorm"

	"gorm.io/driver/postgres"

	"github.com/Raphy42/weekend/core"
)

func postgresDialector() func(dsn string) gorm.Dialector {
	return func(dsn string) gorm.Dialector {
		return postgres.Open(dsn)
	}
}

func init() {
	core.RegisterOnStartHook(func() error {
		globalDriver.Register("postgres", postgresDialector())
	})
	core.RegisterOnStartHook(func() error {
		globalDriver.Register("postgresql", postgresDialector())
	})
}
