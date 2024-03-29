package tramservice

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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

func (m *RPCMessage) PrepReply(in *RPCMessage) {
	m.MessageType = Reply
	m.RequestID = in.RequestID
	m.RPCId = in.RPCId
	m.TranslationId = in.TranslationId
	m.ProcedureID = in.ProcedureID
}

type Tram struct {
	TramID       *uuid.UUID
	CurrentStop  int
	PreviousStop int
}

func (curr *Tram) ToString() string {
	return fmt.Sprintf("%s,%d,%d", curr.TramID.String(), curr.CurrentStop, curr.PreviousStop)
}

func (curr *Tram) FromString(data string) {
	temp := strings.Split(data, ",")
	if len(temp) != 4 {
		panic("Oh oh, couldn't unpack Tram data")
	}
	curr.TramID, _ = uuid.ParseHex(temp[0])
	curr.CurrentStop, _ = strconv.Atoi(strings.TrimSpace(temp[1]))
	curr.PreviousStop, _ = strconv.Atoi(strings.TrimSpace(temp[2]))
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
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		out[i] = fmt.Sprintf("%v", f.Interface())
	}

	output := new(Message)
	output.data = []byte(fmt.Sprintf("%s\n", strings.Join(out, "|")))
	output.length = uint32(len(output.data))
	return output
}

// Unmarshall decodes CsvData from the Message.
// should return new RPCMessage with parsed values.
func (message *Message) Unmarshall() *RPCMessage {
	output := new(RPCMessage)
	tempData := strings.Split(string(message.data), "|")
	if tempData[0] == strconv.Itoa(0) {
		output.MessageType = Request
	} else {
		output.MessageType = Reply
	}
	// a bit ugly but works
	transactionID, _ := strconv.Atoi(tempData[1])
	output.TranslationId = uint32(transactionID)
	output.RPCId, _ = uuid.ParseHex(tempData[2])
	requestID, _ := strconv.Atoi(tempData[3])
	output.RequestID = uint32(requestID)
	procedureID, _ := strconv.Atoi(tempData[4])
	output.ProcedureID = uint32(procedureID)
	output.CsvData = tempData[5]
	status, _ := strconv.Atoi(tempData[6])
	output.Status = uint32(status)
	return output
}

// getRand gets a new delay value for tram simulation
// as per assignment spec.
func genRand() (out int) {
	rand.Seed(time.Now().UTC().UnixNano())
	out = 10 + rand.Intn(10)
	return
}

// debug logging facility for assignment demonstration
func Logger(message string) (err error) {
	file, err := os.OpenFile("tramservice.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(message)
	return
}
