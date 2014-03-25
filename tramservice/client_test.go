package tramservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// init server if hasn't been already
var server Server

func ConnectServer() {
	if server.Status != true {
		go server.init()
	}
}

// Test GetNextStop functionality
func TestGetNextStop(t *testing.T) {
	// init server
	ConnectServer()

	// in/out test input for function
	// consists of tramID, currentstop, previousstop
	tests := [][]int{
		[]int{1, 3, 4, 2},
		[]int{1, 3, 2, 4},
		[]int{1, 1, 2, 2},
		[]int{1, 5, 4, 4},
		[]int{1, 99, 2, -1},
		[]int{1, 99, 2, -1}, // 6th tram on route should break
	}

	testClients := make([]Client, 6)
	for i, b := range testClients {
		b.Init("localhost:1234")
		err := b.RegisterRoute(tests[i][0])
		// error raised if trying to register more than 5 trams
		if i > 5 {
			assert.NotNil(t, err, "RegisterRoute did not fail as expected.")
		}
		b.SetCurrentLocation(tests[i][1], tests[i][2])
		nextStop, err := b.GetNextStop()
		// error raised if invalid current stop
		if tests[i][1] == 99 {
			assert.NotNil(t, err, "NextStop did not fail as expected.")
		} else {
			assert.Nil(t, err, "Error getting next stop.")
		}
		assert.Equal(t, tests[i][3], nextStop, "Next stop wasn't the one expected.")
	}
}

// TestUpdateTramLocation verifies that current Tram locations are updated correctly
// this is done by advancing the tram once from current starting position and checking
// expected result.
func TestUpdateTramLocation(t *testing.T) {
	// init server
	ConnectServer()

	// in/out input for function
	tests := map[int]int{
		96: 24,
	}

	for i, z := range tests {
		b := new(Client)
		b.Init("localhost:1234")
		b.RegisterRoute(i)
		b.AdvanceTram()
		//assert.Nil(t, err, "Error getting next stop.")
		assert.Equal(t, z, b.TramObj.CurrentStop, "Next stop wasn't the one expected.")
	}

}

func TestSequentialPathing(t *testing.T) {
	// init server
	ConnectServer()

	// clear current clients from previous tests
	server.clearClients()

	// list of routes to bind the running trams to
	tests := []int{1, 96, 1}

	// initialise test trams, ready to start moving
	workingClients := make([]Client, len(tests))
	for i, b := range workingClients {
		b.Init("localhost:1234")
		err := b.RegisterRoute(tests[i])
		assert.Nil(t, err, "Error registering route")
	}
	testChannels := make([]chan int, len(tests))
	for i, b := range workingClients {
		go b.AsyncAdvance(testChannels[i])
	}

	select {
	case tram1 := <-testChannels[0]:
		fmt.Println("Tram 1 received")
		if tram1 != 1 {
			go workingClients[0].AsyncAdvance(testChannels[0])
		}
	case tram2 := <-testChannels[1]:
		fmt.Println("Tram 2 received")
		if tram2 != 1 {
			go workingClients[1].AsyncAdvance(testChannels[1])
		}
	case tram3 := <-testChannels[2]:
		fmt.Println("Tram 3 received")
		if tram3 != 2 {
			go workingClients[2].AsyncAdvance(testChannels[2])
		}
	}
}
