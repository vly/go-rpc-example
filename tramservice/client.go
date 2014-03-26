package tramservice

import (
	"errors"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

// Client represents the functional Client
// and includes the open socket, requests count
// and new Tram object for passing to server.
type Client struct {
	socket   *rpc.Client
	requests uint32
	TramObj  *Tram
	routeID  int
}

// Check for errors, and exit if found
func (c *Client) checkError(err error) {
	if err != nil {
		log.Fatalln("Fatal error", err.Error())
		os.Exit(1)
	}
}

// checkID varifies if the incoming RPCMessage has matching IDs
// as the sent ones.
func (c *Client) checkIDs(to *RPCMessage, from *RPCMessage) {
	if to.RPCId.String() != from.RPCId.String() {
		log.Fatalf("Expected %d but received %d\n", to.RPCId, from.RPCId)
	}
	if to.ProcedureID != from.ProcedureID {
		log.Fatalf("Expected %d but received %d\n", to.ProcedureID, from.ProcedureID)
	}
	if to.RequestID != from.RequestID {
		log.Fatalf("Expected %d but received %d\n", to.RequestID, from.RequestID)
	}
	if to.TranslationId != from.TranslationId {
		log.Fatalf("Expected %d but received %d\n", to.TranslationId, from.TranslationId)
	}
}

// checkStatus verifies the response status code.
func (c *Client) checkStatus(response *RPCMessage) error {
	if response.Status != 0 {
		return errors.New("Server responsed with an error")
	}
	return nil
}

// Init initialises Client functionality
// by establishing connection to server and generating
// a new tram object.
func (c *Client) Init(serverIP string) (err error) {
	client, err := rpc.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	c.socket = client

	// generate a new tram

	err = c.genTram()
	if err != nil {
		log.Fatalln("Couldn't gen new tram")
	}
	return
}

// genTram generates a new Tram instance
// including new uuid.
func (c *Client) genTram() (err error) {
	c.TramObj = new(Tram)
	c.TramObj.TramID, err = uuid.NewV4()
	if err != nil {
		log.Fatalln("Error generating UUID")
	}
	return
}

// RegisterRoute enables tram to be bound to a specific route
// this is a prerequisite for issuing any further commands to the server.
func (c *Client) RegisterRoute(routeID int) error {
	// Synchronous call
	c.requests += 1
	//Specify procedureID
	procedureID := uint32(0)
	rpcID, err := uuid.NewV4()
	c.checkError(err)
	data := fmt.Sprintf("%s,%d", c.TramObj.ToString(), routeID)
	c.routeID = routeID
	newMessage := RPCMessage{Request, 1, rpcID, 1, procedureID, data, 0}

	var response *RPCMessage
	var mashalledResponse Message
	marshalledMessage := newMessage.Marshall()
	// direct call: err = c.socket.Call("Server.RegisterTram", &newMessage, &response)
	err = c.socket.Call("Server.CallBroker", &marshalledMessage, &mashalledResponse)
	c.checkError(err)
	response = mashalledResponse.Unmarshall()
	c.checkIDs(&newMessage, response)
	chk := c.checkStatus(response)
	if chk != nil {
		fmt.Println("RegisterRoute: Error registering the tram, route tram limit reached.")
		return chk
	}

	if len(response.CsvData) != 0 {
		data := strings.Split(response.CsvData, ",")
		c.TramObj.CurrentStop, _ = strconv.Atoi(data[0])
		c.TramObj.PreviousStop, _ = strconv.Atoi(data[1])
	}
	return err
}

// AdvanceTram moves the current tram to the next stop
func (c *Client) AdvanceTram() error {
	nextStop, err := c.GetNextStop()
	c.checkError(err)
	Logger(fmt.Sprintf("tram %s received next stop %d", c.TramObj.TramID.String(), nextStop))
	// sleep before executing
	time.Sleep(time.Duration(genRand()) * time.Second)
	err = c.UpdateTramLocation(nextStop)
	return err
}

// GetNextStop requests the next stop ID for the current
// route the tram/client is on. It is directional therefore
// previous stop and current stop IDs are passed as params.
func (c *Client) GetNextStop() (nextStop int, err error) {
	// Synchronous call
	c.requests += 1
	//Specify procedureID
	procedureID := uint32(1)

	rpcID, err := uuid.NewV4()
	c.checkError(err)

	data := fmt.Sprintf("%d,%d,%d", c.routeID, c.TramObj.CurrentStop, c.TramObj.PreviousStop)

	newMessage := RPCMessage{Request, 1, rpcID, 1, procedureID, data, 0}
	marshalledMessage := newMessage.Marshall()

	var response *RPCMessage
	mashalledResponse := new(Message)
	// for direct call: err = c.socket.Call("Server.GetNextStop", &newMessage, &response)

	err = c.socket.Call("Server.CallBroker", &marshalledMessage, &mashalledResponse)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	response = mashalledResponse.Unmarshall()
	c.checkIDs(&newMessage, response)
	chk := c.checkStatus(response)
	if chk != nil {
		fmt.Println("GetNextStop: Error getting next stop, current tram position is not valid.")
		return -1, chk
	}
	nextStop, err = strconv.Atoi(response.CsvData)
	c.checkError(err)
	return
}

// SetCurrentLocation overwrites current stops in local Tram object
func (c *Client) SetCurrentLocation(currentStop int, previousStop int) error {
	c.TramObj.CurrentStop = currentStop
	c.TramObj.PreviousStop = previousStop

	return nil
}

// UpdateTramLocation notifies the server that the tram has arrived at
// the next tram stop.
func (c *Client) UpdateTramLocation(nextStop int) (err error) {
	// increment requests counter
	c.requests += 1
	//Specify procedureID
	procedureID := uint32(2)

	rpcID, err := uuid.NewV4()
	// pass a subset of Tram object to satisfy brief
	csvData := fmt.Sprintf("%s,%d", c.TramObj.TramID.String(), nextStop)
	newMessage := RPCMessage{Request, 1, rpcID, c.requests, procedureID, csvData, 0}

	// carry out the RPC call and process response
	marshalledMessage := newMessage.Marshall()

	var response *RPCMessage
	mashalledResponse := new(Message)
	// direct call: err = c.socket.Call("Server.UpdateTramLocation", &newMessage, response)

	err = c.socket.Call("Server.CallBroker", &marshalledMessage, &mashalledResponse)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	response = mashalledResponse.Unmarshall()
	c.checkIDs(&newMessage, response)
	chk := c.checkStatus(response)
	if chk != nil {
		fmt.Println("UpdateTramLocation: Error updating current location, next stop is not valid.")
		return chk
	}
	// if everything OK, set the nextStop as currentStop
	if len(response.CsvData) == 0 {
		c.TramObj.PreviousStop = c.TramObj.CurrentStop
		c.TramObj.CurrentStop = nextStop
	}

	return
}

// Async wrapper
func (c *Client) AsyncAdvance() {
	c.AdvanceTram()
	return

}
