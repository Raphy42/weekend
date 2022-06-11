package task

const (
	PayloadTypeJSON    = "json"
	PayloadTypeMsgPack = "msgpack"
)

type Manifest struct {
	Name        string
	Payload     []byte
	PayloadType string
}
