package v2

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)

type Client struct {
	socket *rpc.Client
}

func (c *Client) Init() (err error) {
	service := "localhost:1234"
	client, err := rpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	c.socket = client
	return
}

func (c *Client) GetNextStop(data *CurrentLoc) (nextStop int, err error) {
	// Synchronous call
	newMessage := RPCMessage{Request, 1, 3366222, 1, 1, data.ToString(), 1}

	var response RPCMessage
	err = c.socket.Call("Arith.GetNextStop", &newMessage, &response)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Response: %s\n", response.CsvData)
	nextStop, err = strconv.Atoi(response.CsvData)
	return
}

// 	var quot Quotient
// 	err = client.Call("Arith.Divide", args, &quot)

// 	if err != nil {
// 		log.Fatal("arith error:", err)
// 	}
// 	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.TramID, args.CurrentStop, quot.Quo, quot.Rem)
// }
