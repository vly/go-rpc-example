package tramservice

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/rpc"
	"strconv"
	"time"
)

// Client represents the functional Client
// and includes the open socket, requests count
// and new Tram object for passing to server.
type Client struct {
	socket   *rpc.Client
	requests uint32
	tram     *Tram
	routeID  int
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
	c.genTram()
	return
}

// genTram generates a new Tram instance
// including new uuid.
func (c *Client) genTram() (err error) {
	c.tram = new(Tram)
	c.tram.TramID, err = uuid.NewV4()
	if err != nil {
		log.Fatalln("Error generating UUID")
	}
	return
}

func (c *Client) RegisterRoute(routeID int) error {
	// Synchronous call
	c.requests += 1
	rpcID, err := uuid.NewV4()
	data := fmt.Sprintf("%s,%d", c.tram.ToString(), routeID)
	c.routeID = routeID
	newMessage := RPCMessage{Request, 1, rpcID, 1, 1, data, 1}

	var response RPCMessage
	err = c.socket.Call("Server.RegisterTram", &newMessage, &response)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	c.checkIDs(&newMessage, &response)

	fmt.Printf("Response: %s\n", response.CsvData)
	return err
}

// checkID varifies if the incoming RPCMessage has a matching RequestID
// as the sent one.
func (c *Client) checkIDs(to *RPCMessage, from *RPCMessage) {
	if to.RequestID != from.RequestID {
		log.Fatalf("Expected %d but received %d\n", to.RequestID, from.RequestID)
	}
}

// AdvanceTram moves the current tram to the next stop
func (c *Client) AdvanceTram() {
	c.GetNextStop()
	// sleep before executing
	time.Sleep(time.Duration(genRand()) * time.Second)
	c.UpdateTramLocation()
}

// GetNextStop requests the next stop ID for the current
// route the tram/client is on. It is directional therefore
// previous stop and current stop IDs are passed as params.
func (c *Client) GetNextStop() (nextStop int, err error) {
	// Synchronous call
	c.requests += 1
	rpcID, err := uuid.NewV4()
	data := fmt.Sprintf("%d,%d, %d", c.routeID, c.tram.CurrentStop, c.tram.PreviousStop)
	newMessage := RPCMessage{Request, 1, rpcID, 1, 1, data, 1}

	var response RPCMessage
	err = c.socket.Call("Server.GetNextStop", &newMessage, &response)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	c.checkIDs(&newMessage, &response)

	fmt.Printf("Response: %s\n", response.CsvData)
	nextStop, err = strconv.Atoi(response.CsvData)
	return
}

// UpdateTramLocation notifies the server that the tram has arrived at
// the next tram stop.
func (c *Client) UpdateTramLocation() (nextStop int, err error) {
	// increment requests counter
	c.requests += 1
	rpcID, err := uuid.NewV4()
	// pass a subset of Tram object to satisfy brief
	csvData := fmt.Sprintf("%s,%d", c.tram.TramID.String(), c.tram.CurrentStop)
	newMessage := RPCMessage{Request, 1, rpcID, c.requests, 1, csvData, 1}

	// compress the message using a custom marshalling function
	// as described in part two of the assignment
	temp := newMessage.Marshall()
	temp.Unmarshall()

	// carry out the RPC call and process response
	var response RPCMessage
	err = c.socket.Call("Server.UpdateTramLocation", &newMessage, &response)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	c.checkIDs(&newMessage, &response)
	fmt.Printf("Response: %s\n", response.CsvData)
	if len(response.CsvData) != 0 {
		nextStop, err = strconv.Atoi(response.CsvData)
	}
	return
}
