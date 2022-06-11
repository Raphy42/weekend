package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	a := assert.New(t)

	s := NewSlice(1, 2, 3, 4, 5, 6, 7)
	a.Equal(NewSlice(10, 20, 30, 40, 50, 60, 70), Map(s, func(_ int, in int) int {
		return in * 10
	}))
}

func TestFilter(t *testing.T) {
	a := assert.New(t)

	s := NewSlice(1, 2, 3, 4, 5, 6, 7)
	a.Equal(NewSlice(1, 2, 3), Filter(s, func(_, in int) bool {
		return in < 4
	}))
}

func TestUnique(t *testing.T) {
	a := assert.New(t)

	si := NewSlice(1, 2, 3, 3, 4, 5, 5, 6, 7)
	a.Equal(NewSlice(1, 2, 3, 4, 5, 6, 7), Unique(si))
	ss := NewSlice("a", "b", "b")
	a.Equal(NewSlice("a", "b"), Unique(ss))
}

func TestFlatten(t *testing.T) {
	a := assert.New(t)

	s := NewSlice(1, 2, 3)
	a.Equal(NewSlice(1, 2, 3, 1, 2, 3, 1, 2, 3), Flatten(s, s, s))
}
