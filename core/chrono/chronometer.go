package chrono

import (
	"context"
	"time"
)

//Chrono can be used to time things around various parts of the application
type Chrono struct {
	start time.Time
}

//NewChrono creates a `time.Chrono`
func NewChrono() *Chrono {
	return &Chrono{}
}

//Start starts the chronometer normally
func (c *Chrono) Start() {
	c.start = time.Now()
}

//StartContext starts the chronometer and ends it whenever the context is terminated
// returning the duration via a channel
func (c *Chrono) StartContext(ctx context.Context) chan time.Duration {
	c.Start()
	duration := make(chan time.Duration)
	go func() {
		select {
		case <-ctx.Done():
			duration <- time.Now().Sub(c.start)
		}
	}()
	return duration
}

//StartDeferred starts the chronometer and ends it whenever the returned callback is invoked,
// returning the duration via a channel
func (c *Chrono) StartDeferred() (func(), chan time.Duration) {
	root := context.Background()
	child, cancel := context.WithCancel(root)
	return func() {
		cancel()
	}, c.StartContext(child)
}

//Elapsed gets the duration elapsed between the time the chronometer was initialised.
// This method does not reset the internal clock, unless `Chrono.Start` is called.
func (c *Chrono) Elapsed() time.Duration {
	return time.Now().Sub(c.start)
}
