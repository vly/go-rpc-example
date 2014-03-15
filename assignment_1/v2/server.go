package v2

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

type Arith int

func (t *Arith) GetNextStop(in *RPCMessage, out *RPCMessage) error {
	fmt.Println(in.CsvData)
	var data CurrentLoc
	data.FromString(in.CsvData)
	// get array of stops
	stops, err := inDatabase(&data)
	if err != nil {
		return nil
	}
	// check if current stop is in there
	var nextStop int
	for a, b := range stops {
		if b == data.CurrentStop {
			if a != len(stops)-1 {
				nextStop = stops[a+1]
			} else {
				nextStop = stops[a-1]
			}
		}
	}
	out.MessageType = Reply
	out.CsvData = strconv.Itoa(nextStop)
	return nil
}
func (t *Arith) Divide(args *CurrentLoc, quo *Quotient) error {
	if args.CurrentStop == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.TramID / args.CurrentStop
	quo.Rem = args.TramID % args.CurrentStop

	return nil
}

func (t *Arith) checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func (t *Arith) Init() {
	rpc.Register(t)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	t.checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	t.checkError(err)
	/* This works:
	rpc.Accept(listener)
	*/
	/* and so does this:
	 */
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}

// check if tram and stop are in database
func inDatabase(data *CurrentLoc) (stops []int, err error) {
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

// // retrieve next tram stop
// func NextStop(in *RPCMessage, out *RPCMessage) error {
// 	data := strings.Split(in.csvData, ",")
// 	if len(data) != 2 {
// 		return errors.New("not enough params")
// 	}
// 	tram, _ := strconv.Atoi(data[0])
// 	currentStop, _ := strconv.Atoi(data[1])

// 	var nextStop int
// 	if err != nil {
// 		return nil
// 	}
// 	for a, b := range stops {
// 		if b == currentStop {
// 			if a != len(stops)-1 {
// 				nextStop = stops[a+1]
// 			} else {
// 				nextStop = stops[a-1]
// 			}
// 		}
// 	}
// 	out.csvData = strconv.Itoa(nextStop)
// 	return nil
// }
