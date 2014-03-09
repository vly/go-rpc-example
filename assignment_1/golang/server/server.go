package server

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"net"
	"net/rpc"
)

type RPCServer struct{}

func (t *RPCServer) Test(args *string, reply *string) error {
	*reply = "Hello back"
	return nil
}

func Start() {
	var mh codec.MsgpackHandle

	theServer := new(RPCServer)
	rpc.Register(theServer)

	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		rpcCodec := codec.MsgpackSpecRpc.ServerCodec(conn, &mh)
		fmt.Printf("%s\n", "hello")
		go rpc.ServeCodec(rpcCodec)
	}
}
