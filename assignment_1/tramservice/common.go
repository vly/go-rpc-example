package tramservice

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"reflect"
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

// Custom message encoding struct
type Message struct {
	data   []byte
	length uint32
}

// encode CsvData
// should return *Message
func (message *RPCMessage) Marshall() *Message {
	s := reflect.ValueOf(message).Elem()
	out := make([]string, s.NumField())
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		out[i] = fmt.Sprintf("%s:%v", typeOfT.Field(i).Name, f.Interface())
	}

	output := new(Message)
	output.data = []byte(fmt.Sprintf("%s\n", strings.Join(out, "|")))
	output.length = uint32(len(output.data))

	return output
}

// Unmarshall decodes CsvData from the Message.
// should return new RPCMessage with parsed values.
func (message *Message) Unmarshall() (output *RPCMessage) {

	tempData := strings.Split(string(message.data), "|")
	fmt.Println(tempData)
	return output
}
