package main

import (
	"./tramservice"
	"fmt"
)

// Main tramservice server
func main() {
	fmt.Println("Starting Tramservice server")
	server := new(tramservice.Server)
	server.Init()
}
