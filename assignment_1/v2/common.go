package v2

import (
	"fmt"
	"strconv"
	"strings"
)

type MessageType int

const (
	Request MessageType = iota
	Reply
)

type RPCMessage struct {
	MessageType   MessageType
	TranslationId uint32
	RPCId         uint32
	RequestID     uint32
	ProcedureID   uint32
	CsvData       string
	Status        uint32
}

type CurrentLoc struct {
	TramID      int
	CurrentStop int
}

type Quotient struct {
	Quo, Rem int
}

func (curr *CurrentLoc) ToString() string {
	return fmt.Sprintf("%d,%d", curr.TramID, curr.CurrentStop)
}

func (curr *CurrentLoc) FromString(data string) {
	temp := strings.Split(data, ",")
	if len(temp) != 2 {
		panic("Oh oh, couldn't unpack")
	}
	curr.TramID, _ = strconv.Atoi(temp[0])
	curr.CurrentStop, _ = strconv.Atoi(temp[1])
}
