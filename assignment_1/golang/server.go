package main

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"net"
	"net/rpc"
)

type RPCServer struct{}

type RPCMessage struct {
	 enum {Request, Reply} messageType; /* same size as an unsigned int */ 
	 unsigned int TransactionId /* transaction id */ 
	 unsigned int RPCId; /* Globally unique identifier */ 
	 unsigned int RequestId; /* Client request message counter */ 
	 unsigned int procedureId; /* e.g.(1,2,3,4) */ 
	 char csv_data[] /* data as comma separated values*/ 
	 unsigned int length; /*length of data in cvs_data*/ 
	 unsigned int status; /*status of the transaction*/ 
}

func (t *RPCServer) Test(args *string, reply *string) error {
	*reply = "Hello back"
	return nil
}

func main() {
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
