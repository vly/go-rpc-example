package v2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test GetNextStop functionality
func TestGetNextStop(t *testing.T) {
	// init server
	server := new(Arith)
	go server.Init()
	// in/out test input for function
	tests := map[CurrentLoc]int{
		CurrentLoc{1, 5}:    4,
		CurrentLoc{1, 1}:    2,
		CurrentLoc{101, 34}: 5,
	}

	// initialise client
	client := new(Client)
	err := client.Init()
	assert.Nil(t, err, "Error initialising client")

	// prep data and get reply
	for a, b := range tests {
		result, err := client.GetNextStop(&a)
		assert.Nil(t, err, "Error getting next stop")
		fmt.Printf("%d\n", result)
		assert.Equal(t, b, result, "Received unexpected results")
	}
}
