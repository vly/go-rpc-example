package server

import (
	"github.com/golang/glog"
	//"github.com/ugorji/go/codec"
	"log"
	"net"
)

const (
	SERVER_PORT = "4444"
)

type Client struct {
	uid      string
	address  string
	messages []RPCMessage
}

type ServerState struct {
	status         bool
	socket         *net.UDPConn
	currentClients int64
	clients        []Client
}

func (s *ServerState) Bind() (err error) {
	glog.Info("Initialising server listener")
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

func (s *ServerState) Listen() {
	buf := make([]byte, 1024)
	glog.Infoln("Starting to read")
	for {
		dsize, _, err := s.socket.ReadFromUDP(buf)
		if err != nil {
			log.Panic(err)
		}
		go ServerProcessData(buf[:dsize])
	}

}

// Processes incoming datagrams, unmarshalls and stores for pickup
func ServerProcessData(data []byte) {
	log.Print(string(data))
}
