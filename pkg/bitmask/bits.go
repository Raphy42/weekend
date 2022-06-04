package bitmask

import (
	"github.com/Raphy42/weekend/pkg/kind"
)

func Set[B kind.Bit](b, flag B) B {
	return b | flag
}

func Clear[B kind.Bit](b, flag B) B {
	return b &^ flag
}

func Toggle[B kind.Bit](b, flag B) B {
	return b ^ flag
}

func Has[B kind.Bit](b, flag B) bool {
	return b&flag != 0
}
