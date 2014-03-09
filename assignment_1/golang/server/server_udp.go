package server

import (
	"github.com/golang/glog"
	"log"
	"net"
	"time"
)

func ProtoClient(message string) error {
	time.Sleep(1)
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:4444")
	con, _ := net.DialUDP("udp", nil, serverAddr)
	defer con.Close()
	_, err := con.Write([]byte(message))
	return err
}

const (
	SERVER_PORT = "4444"
)

func ServerBind() (sock *net.UDPConn, err error) {
	glog.Info("Initialising server listener")
	addr, err := net.ResolveUDPAddr("udp", ":"+SERVER_PORT)
	if err != nil {
		glog.Fatalln("Couldn't resolve local address")
	}
	sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		glog.Fatalln("Failed to bind")
	}

	return
}

func ServerListen(sock *net.UDPConn) {
	buf := make([]byte, 1024)
	glog.Infoln("Starting to read")
	go ProtoClient("testing")
	for {
		dsize, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			log.Panic(err)
		}
		log.Print(string(buf[:dsize]))
	}

}
