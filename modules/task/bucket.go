package task

import (
	"sync"
)

type Bucket struct {
	lock           sync.RWMutex
	maxConcurrency int32
	inflight       int32
}
type SlotReleaser func()

func releaseSlotNoOp() {}

func newBucket(max int32) *Bucket {
	return &Bucket{
		maxConcurrency: max,
		inflight:       0,
	}
}

func (b *Bucket) HasSlot() bool {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.inflight < b.maxConcurrency
}

func (b *Bucket) releaseSlot() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.inflight--
}

func (b *Bucket) Reserve() (bool, SlotReleaser) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.inflight > b.maxConcurrency {
		return false, releaseSlotNoOp
	}

	b.inflight += 1
	return true, b.releaseSlot
}
