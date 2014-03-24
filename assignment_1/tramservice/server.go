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

// const (
// 	// Default server address to bind to
// 	ServerAddress = ":1234"
// )

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

// Initialise Server, registering RPC service and binding to port 1234
func (t *Server) Init() {

	t.Status = true
	rpc.Register(t)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
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
			return errors.New("Maximum trams reached on this route.")
		}
	}
	t.Routes[routeID] = append(t.Routes[routeID], data)
	t.Clients[data.TramID.String()] = &Record{routeID, data}
	//t.Clients[data.] .PushFront([]interface{data)
	return nil
}

// CallBroker routes calls to various functions
// this is not really necessary as I am exposing selected methods
// through RPC registry of Server struct but... for the sake of assignment
// here it is.
func (t *Server) CallBroker(in *RPCMessage, out *RPCMessage) error {
	switch in.ProcedureID {
	case 0:
		return t.RegisterTram(in, out)
	case 1:
		return t.GetNextStop(in, out)
	case 2:
		return t.UpdateTramLocation(in, out)
	}
	out.PrepReply(in)
	out.Status = 1
	return nil
}

// RegisterTram functionality
// enables trams to be attached to specific routes.
func (t *Server) RegisterTram(in *RPCMessage, out *RPCMessage) error {
	glog.Infoln("RegisterTram received: " + in.CsvData)
	out.PrepReply(in)
	tempSplit := strings.Split(in.CsvData, ",")
	routeID, err := strconv.Atoi(tempSplit[len(tempSplit)-1])
	if err != nil {
		glog.Fatalln("Error splitting out tram route from RPCMessage data.")
	}
	stops, err := inDatabase(routeID)
	if err != nil {
		glog.Fatalln("Route doesn't exist")
	}

	var data Tram
	data.FromString(in.CsvData)
	err = t.addClient(&data, routeID)
	if err != nil {
		out.Status = 1
	} else {
		// pass current and previous stops to client
		// these represent the starting (depo) location
		out.CsvData = fmt.Sprintf("%d,%d", stops[0], stops[1])
	}

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
	return errors.New("Tram not registered")
}

// getStats prints out list of current clients
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
	//glog.Infoln("Getnextstop received: " + in.CsvData)

	out.PrepReply(in)

	// get array of stops
	tempSplit := strings.Split(in.CsvData, ",")

	routeID, _ := strconv.Atoi(tempSplit[0])
	currentStop, err := strconv.Atoi(strings.TrimSpace(tempSplit[1]))
	if err != nil {
		t.checkError(err)
	}
	previousStop, _ := strconv.Atoi(tempSplit[2])
	stops, err := inDatabase(routeID)
	// if tramID is not present, return -1
	if err != nil {
		out.Status = 1
		out.CsvData = "-1"
		return nil
	}
	// check if current stop is in there and find the previous stop
	var nextStop int = -2
	for a, b := range stops {
		if b == currentStop {
			if a != len(stops)-1 {
				if previousStop == stops[a+1] && a != 0 {
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
		out.Status = 1
	}

	out.CsvData = strconv.Itoa(nextStop)
	return nil
}

// UpdateTramLocation function
// returns empty CSV if OK, otherwise -1
func (t *Server) UpdateTramLocation(in *RPCMessage, out *RPCMessage) error {
	glog.Infoln("Updatetramlocation received: " + in.CsvData)
	out.PrepReply(in)
	// process input params
	tempSplit := strings.Split(in.CsvData, ",")
	if len(tempSplit) != 2 {
		glog.Fatalln("Incorrect number of params supplied.")
	}
	newStop, _ := strconv.Atoi(tempSplit[1])

	// check if tram has been registered
	if _, ok := t.Clients[tempSplit[0]]; ok {
		t.Clients[tempSplit[0]].Data.PreviousStop = t.Clients[tempSplit[0]].Data.PreviousStop
		t.Clients[tempSplit[0]].Data.CurrentStop = newStop
		return nil
	}

	glog.Infoln("Tram not registered.")
	out.CsvData = "-1"
	out.Status = 1
	return nil

}
