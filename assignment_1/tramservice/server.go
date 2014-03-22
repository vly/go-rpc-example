/*
	Package tramservice provides the server and client functionality
	for tracking trams registered on the network.
*/

package tramservice

import (
	//"container/list"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

// Main server structure
type Server struct {
	Clients map[string]*Record
	Routes  map[int][]*Tram
	Status  bool
}

type Record struct {
	RouteID int
	Data    *Tram
}

const (
	// Default server address to bind to
	ServerAddress = ":1234"
)

// Initialise Server, registering RPC service and binding to port 1234
func (t *Server) Init() {

	t.Status = true
	rpc.Register(t)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ServerAddress)
	t.checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	t.checkError(err)

	// create new map for tracking routes
	t.Routes = make(map[int][]*Tram)
	t.Clients = make(map[string]*Record)

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
		glog.Infoln("Fatal error", err.Error())
		os.Exit(1)
	}
}

// Make sure that a tram route doesn't exceed 5 trams
func (t *Server) addClient(data *Tram, routeID int) error {
	if _, ok := t.Routes[routeID]; ok {
		if len(t.Routes[routeID]) == MAX_ROUTE_TRAM {
			return nil
		}
	}
	t.Routes[routeID] = append(t.Routes[routeID], data)
	t.Clients[data.TramID.String()] = &Record{routeID, data}
	//t.Clients[data.] .PushFront([]interface{data)
	return nil
}

func (t *Server) checkRoute(data *Tram) (int, error) {
	return 1, nil
}

// RegisterTram functionality
// enables trams to be attached to specific routes.
func (t *Server) RegisterTram(in *RPCMessage, out *RPCMessage) error {
	glog.Infoln("RegisterTram received: " + in.CsvData)
	tempSplit := strings.Split(in.CsvData, ",")
	routeID, err := strconv.Atoi(tempSplit[len(tempSplit)-1])
	if err != nil {
		glog.Fatalln("Error splitting out tram route")
	}
	var data Tram
	data.FromString(in.CsvData)
	err = t.addClient(&data, routeID)
	return nil
}

// UnregisterTram removes tramID from current Clients list
// and the route register
func (t *Server) UnregisterTram(in *RPCMessage, out *RPCMessage) error {
	return nil
}

// update Clients list
func (t *Server) updateClient(data *Tram) error {
	if _, ok := t.Clients[data.TramID.String()]; ok {
		t.Clients[data.TramID.String()].Data = data
		return nil
	}
	// not sure if this is required, maybe able to traverse list.List even if empty
	// if t.Clients.Len() == 0 {
	// 	//t.addClient(data)

	// }
	// check if tramID already exists, and if it does update the currentStop
	// for e := t.Clients.Front(); e != nil; e = e.Next() {
	// 	if e.Value.(*Tram).TramID == data.TramID {
	// 		e.Value = data
	// 		return nil
	// 	}
	// }

	return errors.New("Tram not registered")
}

func (t *Server) getStats() {
	fmt.Printf("Current clients: %d\n", len(t.Clients))
}

// inDatabase checks whether a route is valid
func inDatabase(routeID int) (stops []int, err error) {
	ROUTES := map[int][]int{
		1:   {1, 2, 3, 4, 5},
		96:  {23, 24, 2, 34, 22},
		101: {123, 11, 22, 34, 5, 4, 7},
		109: {88, 87, 85, 80, 9, 7, 2, 1},
		112: {110, 123, 11, 22, 34, 33, 29, 4},
	}
	if values, ok := ROUTES[routeID]; ok {
		stops = values
	} else {
		err = errors.New("No such tram route found.")
	}
	return
}

// GetNextStop functionality
// is directional
func (t *Server) GetNextStop(in *RPCMessage, out *RPCMessage) error {
	glog.Infoln("Getnextstop received: " + in.CsvData)
	var data Tram
	data.FromString(in.CsvData)
	out.MessageType = Reply
	out.RequestID = in.RequestID
	// get array of stops
	tempSplit := strings.Split(in.CsvData, ",")
	routeID, _ := strconv.Atoi(tempSplit[len(tempSplit)-1])
	stops, err := inDatabase(routeID)

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
	glog.Infoln("Updatetramlocation received: " + in.CsvData)
	out.MessageType = Reply
	out.RequestID = in.RequestID

	// process input params
	tempSplit := strings.Split(in.CsvData, ",")
	if len(tempSplit) != 2 {
		glog.Fatalln("Incorrect number of params supplied.")
	}
	newStop, _ := strconv.Atoi(tempSplit[1])

	var data Tram
	data.FromString(in.CsvData)

	// check if tram has been registered
	if _, ok := t.Clients[tempSplit[0]]; ok {
		t.Clients[data.TramID.String()].Data.PreviousStop = t.Clients[data.TramID.String()].Data.PreviousStop
		t.Clients[data.TramID.String()].Data.CurrentStop = newStop
		return nil
	}

	glog.Infoln("Tram not registered.")
	out.CsvData = "-1"
	return nil

	// // get array of stops
	// stops, err := inDatabase(&data)

	// // if tramID is not present, return -1
	// if err != nil {
	// 	out.CsvData = "-1"
	// 	return nil
	// }
	// if data.CurrentStop != 0 {
	// 	// check if current stop is in there and find the previous stop
	// 	for _, b := range stops {
	// 		if b == data.CurrentStop {
	// 			t.updateClient(&data)
	// 			return nil
	// 		}
	// 	}
	// } else {
	// 	glog.Infoln("Initial request")
	// 	data.CurrentStop = stops[0]
	// 	data.PreviousStop = stops[1]
	// 	glog.Infoln(data.CurrentStop)
	// 	t.updateClient(&data)
	// 	return nil
	//}

}
