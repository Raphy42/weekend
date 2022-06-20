package gorm

import (
	"net/url"

	"gorm.io/gorm"

	"github.com/Raphy42/weekend/pkg/concurrent_set"
)

type Driver struct {
	registry concurrent_set.Set[string, *func(dsn string) gorm.Dialector]
}

func newDriver() *Driver {
	return &Driver{
		registry: concurrent_set.New[string, *func(dsn string) gorm.Dialector](),
	}
}

var (
	globalDriver = newDriver()
)

func (d *Driver) Register(name string, dialectFn func(dsn string) gorm.Dialector) {
	d.registry.Insert(name, &dialectFn)
}

func (d *Driver) Dialect(dsn url.URL) (*gorm.Dialector, bool) {
	dialect, ok := d.registry.Get(dsn.Scheme)
	if !ok {
		return nil, false
	}
	dialector := (*dialect)(dsn.String())
	return &dialector, true
}
