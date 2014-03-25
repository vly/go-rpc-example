package main

import (
	"fmt"
	"github.com/vly/COSC1170/tramservice"
)

// Main tramservice server
func main() {
	fmt.Println("Starting Tramservice server")
	server := new(tramservice.Server)
	server.Init()
}
