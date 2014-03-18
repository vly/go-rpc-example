package main

import (
	"./tramservice"
	"fmt"
)

// start tramservice client
func main() {
	fmt.Println("Starting Tramservice client")
	client := new(tramservice.Client)
	err := client.Init("localhost:1234")

	if err != nil {
		fmt.Println("Error connecting to server.")
	}
}
