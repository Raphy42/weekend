package slice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	a := assert.New(t)

	s := New(1, 2, 3, 4, 5, 6, 7)
	a.Equal(New(10, 20, 30, 40, 50, 60, 70), Map(s, func(_ int, in int) int {
		return in * 10
	}))
}

func TestFilter(t *testing.T) {
	a := assert.New(t)

	s := New(1, 2, 3, 4, 5, 6, 7)
	a.Equal(New(1, 2, 3), Filter(s, func(_, in int) bool {
		return in < 4
	}))
}

func TestUnique(t *testing.T) {
	a := assert.New(t)

	si := New(1, 2, 3, 3, 4, 5, 5, 6, 7)
	a.Equal(New(1, 2, 3, 4, 5, 6, 7), Unique(si))
	ss := New("a", "b", "b")
	a.Equal(New("a", "b"), Unique(ss))
}

func TestFlatten(t *testing.T) {
	a := assert.New(t)

	s := New(1, 2, 3)
	a.Equal(New(1, 2, 3, 1, 2, 3, 1, 2, 3), Flatten(s, s, s))
}

func TestContains(t *testing.T) {
	a := assert.New(t)

	a.True(Contains(New("a", "b", "c"), "a"))
	a.False(Contains(New(1, 2, 3, 4), 10))
	now := time.Now()
	a.True(Contains(New(time.Now(), now, time.Now()), now))
}
