package gorm

import (
	"net/url"
	"sync"

	"gorm.io/gorm"
)

type Driver struct {
	sync.RWMutex
	registry map[string]func(dsn string) gorm.Dialector
}

func newDriver() *Driver {
	return &Driver{
		registry: make(map[string]func(dsn string) gorm.Dialector),
	}
}

var (
	globalDriver = newDriver()
)

func (d *Driver) Register(name string, dialectFn func(dsn string) gorm.Dialector) {
	d.Lock()
	defer d.Unlock()

	d.registry[name] = dialectFn
}

func (d *Driver) Dialect(dsn url.URL) (*gorm.Dialector, bool) {
	d.RLock()
	defer d.RUnlock()

	dialect, ok := d.registry[dsn.Scheme]
	if !ok {
		return nil, false
	}
	dialector := dialect(dsn.String())
	return &dialector, true
}
