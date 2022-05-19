package bitmask

func Set(b, flag int16) int16 {
	return b | flag
}

func Clear(b, flag int16) int16 {
	return b &^ flag
}

func Toggle(b, flag int16) int16 {
	return b ^ flag
}

func Has(b, flag int16) bool {
	return b&flag != 0
}
