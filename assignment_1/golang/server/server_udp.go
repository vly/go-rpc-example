package server

import (
	"errors"
	"github.com/golang/glog"
	"github.com/ugorji/go/codec"
	"io"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"strconv"
)

const (
	SERVER_PORT = "4444"
)

type TramServer struct {
	status         bool
	socket         *net.UDPConn
	currentClients int64
	clients        map[int]int
	ch             chan string
	fn             *SharedFunctions
}

type SharedFunctions int

// bind server to local port
func (s *TramServer) Bind() (err error) {
	glog.Info("Initialising server listener")
	s.ch = make(chan string)

	addr, err := net.ResolveUDPAddr("udp", ":"+SERVER_PORT)
	if err != nil {
		glog.Fatalln("Couldn't resolve local address")
	}
	s.socket, err = net.ListenUDP("udp", addr)

	if err != nil {
		glog.Fatalln("Failed to bind")
	}

	return
}

// Initialise listening for datagrams
func (s *TramServer) Listen() {
	buf := make([]byte, 1024)
	exp := new(SharedFunctions)
	s.fn = exp
	rpc.Register(exp)
	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
	var (
		//r io.Reader

		h = &mh // or mh to use msgpack
	)
	glog.Infoln("Starting to read")
	for {
		rpcCodec := codec.GoRpc.ServerCodec(s.socket, h)
		rpc.ServeCodec(rpcCodec)
		dsize, addr, err := s.socket.ReadFromUDP(buf)
		if err != nil {
			log.Panic(err)
		}
		go s.ServerProcessData(buf[:dsize])
		answer := <-s.ch
		if len(answer) != 0 {
			log.Println("Received: " + answer + " from " + addr.IP.String() + ":" + strconv.Itoa(addr.Port))
			break
		}
	}

}

// Processes incoming datagrams, unmarshalls and stores for pickup
func (s *TramServer) ServerProcessData(data []byte) {
	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
	var (
		r io.Reader
		// w io.Writer
		b []byte
		h = &mh // or mh to use msgpack
	)

	dec := codec.NewDecoder(r, h)
	dec = codec.NewDecoderBytes(b, h)
	err := dec.Decode(data)
	log.Println(data)
	if err != nil {
		log.Println("error decoding " + err.Error())
	}
	s.ch <- string(data)
}

// check if tram and stop are in database
func (fn *SharedFunctions) inDatabase(tram int, currentStop int) (stops []int, err error) {
	ROUTES := map[int][]int{
		1:   {1, 2, 3, 4, 5},
		96:  {23, 24, 2, 34, 22},
		101: {123, 11, 22, 34, 5, 4, 7},
		109: {88, 87, 85, 80, 9, 7, 2, 1},
		112: {110, 123, 11, 22, 34, 33, 29, 4},
	}
	if values, ok := ROUTES[tram]; ok {
		stops = values
	} else {
		err = errors.New("No such tram route found.")
	}
	return
}

// retrieve next tram stop
func (fn *SharedFunctions) RetrieveNextStop(tram int, currentStop int) (nextStop int, err error) {
	stops, err := fn.inDatabase(tram, currentStop)
	if err != nil {
		return
	}
	for a, b := range stops {
		if b == currentStop {
			if a != len(stops)-1 {
				nextStop = stops[a+1]
			} else {
				nextStop = stops[a-1]
			}
		}
	}

	return
}

// update current tram location
func (fn *SharedFunctions) UpdateTramLocation(tram int, currentStop int) (err error) {
	stops, err := fn.inDatabase(tram, currentStop)
	if err != nil {
		return
	}
	for _, b := range stops {
		if b == currentStop {
			// if _, ok := fn.clients[tram]; ok {
			// 	//s.clients[tram] = currentStop
			log.Println("yep")
			// 	return
			// } else {
			// 	err = errors.New("Looks like an unregistered client")
			// 	return
			// }

		}
	}
	err = errors.New("No such stop found")
	return
}
