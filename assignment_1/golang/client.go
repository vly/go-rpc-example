package client

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"net"
	"net/rpc"
)

func main() {
	var mh codec.MsgpackHandle

	conn, err := net.Dial("tcp", "localhost:5555")
	if err != nil {
		panic(err)
	}

	rpcCodec := codec.MsgpackSpecRpc.ClientCodec(conn, &mh)
	client := rpc.NewClientWithCodec(rpcCodec)

	var reply string
	err = client.Call("TestServer.Test", "", &reply)
	if err != nil {
		panic(err)
	}

	fmt.Println(reply)
	// Output: reply
}
