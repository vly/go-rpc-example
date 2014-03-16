package v2

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

// Main server structure
type Server struct {
	Clients list.List
	Status  bool
}

// Initialise Server, registering RPC service and binding to port 1234
func (t *Server) init() {
	t.Status = true
	rpc.Register(t)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	t.checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	t.checkError(err)

	rpc.Accept(listener)

	// case want to add additional socket accounting
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		continue
	// 	}
	// 	rpc.ServeConn(conn)
	// }
}

// Check for errors, and exit if found
func (t *Server) checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// update Clients list
func (t *Server) updateClient(data *Tram) {
	if t.Clients.Len() == 0 {
		t.Clients.PushFront(data)
		return
	}
	// check if tramID already exists, and if it does update the currentStop
	for e := t.Clients.Front(); e != nil; e = e.Next() {
		if e.Value.(*Tram).TramID == data.TramID {
			e.Value = data
			return
		}
	}

	// otherwise add to front of the list
	t.Clients.PushFront(data)
}

func (t *Server) getStats() {
	fmt.Printf("Current clients: %d\n", t.Clients.Len())
}

// check if tram and stop are in database
func inDatabase(data *Tram) (stops []int, err error) {
	ROUTES := map[int][]int{
		1:   {1, 2, 3, 4, 5},
		96:  {23, 24, 2, 34, 22},
		101: {123, 11, 22, 34, 5, 4, 7},
		109: {88, 87, 85, 80, 9, 7, 2, 1},
		112: {110, 123, 11, 22, 34, 33, 29, 4},
	}
	if values, ok := ROUTES[data.TramID]; ok {
		stops = values
	} else {
		err = errors.New("No such tram route found.")
	}
	return
}

// GetNextStop functionality
// is directional
func (t *Server) GetNextStop(in *RPCMessage, out *RPCMessage) error {
	fmt.Println("Getnextstop received: " + in.CsvData)
	var data Tram
	data.FromString(in.CsvData)
	out.MessageType = Reply
	out.RequestID = in.RequestID
	// get array of stops
	stops, err := inDatabase(&data)

	// if tramID is not present, return -1
	if err != nil {
		out.CsvData = "-1"
		return nil
	}
	// check if current stop is in there and find the previous stop
	var nextStop int = -2
	for a, b := range stops {
		if b == data.CurrentStop {
			if a != len(stops)-1 {
				if data.PreviousStop == stops[a+1] && a != 0 {
					nextStop = stops[a-1]
				} else {
					nextStop = stops[a+1]
				}
			} else {
				nextStop = stops[a-1]
			}
		}
	}
	if nextStop == -2 {
		nextStop = -1
	}
	out.CsvData = strconv.Itoa(nextStop)
	return nil
}

// UpdateTramLocation function
// returns empty CSV if OK, otherwise -1
func (t *Server) UpdateTramLocation(in *RPCMessage, out *RPCMessage) error {
	fmt.Println("Updatetramlocation received: " + in.CsvData)
	var data Tram
	data.FromString(in.CsvData)

	out.MessageType = Reply
	out.RequestID = in.RequestID
	// get array of stops
	stops, err := inDatabase(&data)

	// if tramID is not present, return -1
	if err != nil {
		out.CsvData = "-1"
		return nil
	}
	// check if current stop is in there and find the previous stop
	for _, b := range stops {
		if b == data.CurrentStop {
			t.updateClient(&data)
			return nil
		}
	}

	out.CsvData = "-1"
	return nil
}
