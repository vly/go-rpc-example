package main

import (
	"./tramservice"
	"fmt"
)

const (
	// Default server address
	ServerAddress = "localhost:1234"
)

// Main tramservice client
// registers on a specific route
// then unregisters on exit
func main() {
	fmt.Println("Starting Tramservice client")
	client := new(tramservice.Client)
	err := client.Init(ServerAddress)

	if err != nil {
		fmt.Println("Error connecting to server.")
	}
}
