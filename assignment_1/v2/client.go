package v2

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)

type Client struct {
	socket   *rpc.Client
	requests uint32
}

func (c *Client) Init(serverIP string) (err error) {
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
	newMessage := RPCMessage{Request, 1, 3366222, 1, 1, data.ToString(), 1}

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
	newMessage := RPCMessage{Request, 1, 3366222, c.requests, 1, data.ToString(), 1}

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
