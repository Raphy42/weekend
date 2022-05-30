package chrono

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeDuration(t *testing.T) {
	a := assert.New(t)

	a.Equal("2m", HumanizeDuration(time.Minute*2))
}