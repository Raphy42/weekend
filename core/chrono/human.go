package chrono

import (
	"fmt"
	"math"
	"time"
)

const (
	Day = time.Hour * 24
)

var (
	formatSlice = []string{"ns", "µs", "ms", "s", "m", "h", "d"}
)

func humanizeDuration(rem time.Duration, acc string) string {
	if rem == 0 {
		return acc
	}
	digits := int(math.Ceil(math.Log10(float64(rem.Nanoseconds()))))
	for idx, format := range formatSlice {
		offset := idx + 1
		if digits%offset == 0 {
			value := math.Pow(10.0, float64(offset))
			significant := math.Floor(float64(rem.Nanoseconds()) / value)
			rem -= time.Duration(significant * value)
			return humanizeDuration(rem, acc+fmt.Sprintf("%d%s", int64(significant), format))
		}
	}
	return acc
}

func HumanizeDuration(duration time.Duration) string {
	return humanizeDuration(duration, "")
}
