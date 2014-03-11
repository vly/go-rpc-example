package server

type MessageType int

const (
	Request MessageType = iota
	Reply
)

type RPCMessage struct {
	messageType   MessageType
	translationId uint32
	RPCId         uint32
	requestID     uint32
	procedureID   uint32
	csvData       string
	status        uint32
}

func (r *RPCMessage) Marshall() []byte {
	return []byte("yep")
}
