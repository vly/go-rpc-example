package v2

import (
	"fmt"
	"strconv"
	"strings"
)

// define system limits
const (
	MAX_ROUTE_TRAM  int = 5
	MIN_ROUTE_STOPS int = 5
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

type Tram struct {
	TramID       int
	CurrentStop  int
	PreviousStop int
}

type Quotient struct {
	Quo, Rem int
}

func (curr *Tram) ToString() string {
	return fmt.Sprintf("%d,%d,%d", curr.TramID, curr.CurrentStop, curr.PreviousStop)
}

func (curr *Tram) FromString(data string) {
	temp := strings.Split(data, ",")
	if len(temp) != 3 {
		panic("Oh oh, couldn't unpack")
	}
	curr.TramID, _ = strconv.Atoi(temp[0])
	curr.CurrentStop, _ = strconv.Atoi(temp[1])
	curr.PreviousStop, _ = strconv.Atoi(temp[2])
}
