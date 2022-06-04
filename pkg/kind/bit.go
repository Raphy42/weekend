package kind

type Bit16 interface {
	uint8 | int16
}

type Bit32 interface {
	uint16 | int32
}

type Bit64 interface {
	uint32 | int64
}

type Bit128 interface {
	uint64
}

type Bit interface {
	Bit16 | Bit32 | Bit64 | Bit128
}
