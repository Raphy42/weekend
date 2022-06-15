package gorm

import (
	"sync"

	"github.com/palantir/stacktrace"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

type DB struct {
	sync.RWMutex
	db *gorm.DB
}

func newDB(dialect *gorm.Dialector) (*DB, error) {
	db, err := gorm.Open(*dialect, &gorm.Config{})
	if err != nil {
		return nil, stacktrace.Propagate(err, "could not build gorm DB instance")
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, stacktrace.Propagate(err, "could not install gorm opentelemetry plugin")
	}

	return &DB{
		db: db,
	}, nil
}
