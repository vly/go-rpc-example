package tramservice

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
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
	RPCId         *uuid.UUID
	RequestID     uint32
	ProcedureID   uint32
	CsvData       string
	Status        uint32
}

type Tram struct {
	TramID       *uuid.UUID
	Route        int
	CurrentStop  int
	PreviousStop int
}

type Quotient struct {
	Quo, Rem int
}

func (curr *Tram) ToString() string {
	return fmt.Sprintf("%s,%d,%d,%d", curr.TramID.String(), curr.Route, curr.CurrentStop, curr.PreviousStop)
}

func (curr *Tram) FromString(data string) {
	temp := strings.Split(data, ",")
	if len(temp) != 4 {
		panic("Oh oh, couldn't unpack")
	}
	curr.TramID, _ = uuid.ParseHex(temp[0])
	curr.Route, _ = strconv.Atoi(temp[1])
	curr.CurrentStop, _ = strconv.Atoi(temp[2])
	curr.PreviousStop, _ = strconv.Atoi(temp[3])
}
