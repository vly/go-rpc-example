package server

import (
	"github.com/golang/glog"
	//"github.com/ugorji/go/codec"
	"errors"
	"log"
	"net"
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
}

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

func (s *TramServer) Listen() {
	buf := make([]byte, 1024)
	glog.Infoln("Starting to read")
	for {
		dsize, _, err := s.socket.ReadFromUDP(buf)
		if err != nil {
			log.Panic(err)
		}
		go s.ServerProcessData(buf[:dsize])
		answer := <-s.ch
		if len(answer) != 0 {
			log.Println(answer)
			break
		}
	}

}

// Processes incoming datagrams, unmarshalls and stores for pickup
func (s *TramServer) ServerProcessData(data []byte) {
	s.ch <- string(data)
}

func (s *TramServer) inDatabase(tram int, currentStop int) (stops []int, err error) {
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

func (s *TramServer) RetrieveNextStop(tram int, currentStop int) (nextStop int, err error) {
	stops, err := s.inDatabase(tram, currentStop)
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

func (s *TramServer) UpdateTramLocation(tram int, currentStop int) (err error) {
	stops, err := s.inDatabase(tram, currentStop)
	if err != nil {
		return
	}
	for _, b := range stops {
		if b == currentStop {
			if _, ok := s.clients[tram]; ok {
				s.clients[tram] = currentStop
				return
			} else {
				err = errors.New("Looks like an unregistered client")
				return
			}
		}
	}
	err = errors.New("No such stop found")
	return
}
