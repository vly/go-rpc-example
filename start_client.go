package main

import (
	"fmt"
	"github.com/vly/COSC1170/tramservice"
)

const (
	// Default server address
	ServerAddress = "localhost:1234"
)

// Main tramservice client
// registers on a specific route
// then unregisters on exit
func mainClient() {
	fmt.Println("Starting Tramservice client")
	client := new(tramservice.Client)
	err := client.Init(ServerAddress)

	if err != nil {
		fmt.Println("Error connecting to server.")
	}
}
