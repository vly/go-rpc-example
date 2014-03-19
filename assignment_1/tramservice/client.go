package tramservice

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"math/rand"
	"net/rpc"
	"strconv"
	"time"
)

type Client struct {
	socket   *rpc.Client
	requests uint32
	tramID   *uuid.UUID
}

func genRand() (out int) {
	rand.Seed(time.Now().UTC().UnixNano())
	out = 10 + rand.Intn(10)
	return
}

func (c *Client) Init(serverIP string) (err error) {
	c.tramID, err = uuid.NewV4()
	if err != nil {
		log.Fatalln("Error generating UUID")
	}
	client, err := rpc.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	c.socket = client
	return
}

func (c *Client) checkIDs(to *RPCMessage, from *RPCMessage) {
	if to.RequestID != from.RequestID {
		log.Fatalf("Expected %d but received %d\n", to.RequestID, from.RequestID)
	}
}

func (c *Client) GetNextStop(data *Tram) (nextStop int, err error) {
	// Synchronous call
	c.requests += 1
	rpcID, _ := uuid.NewV4()
	newMessage := RPCMessage{Request, 1, rpcID, 1, 1, data.ToString(), 1}

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

func (c *Client) UpdateTramLocation(data *Tram) (nextStop int, err error) {
	c.requests += 1
	rpcID, _ := uuid.NewV4()
	newMessage := RPCMessage{Request, 1, rpcID, c.requests, 1, data.ToString(), 1}
	temp := newMessage.Marshall()
	temp.Unmarshall()

	// sleep before executing
	time.Sleep(time.Duration(genRand()) * time.Second)

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
