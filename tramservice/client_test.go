package tramservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	Logger("server initiated")

	// clear current clients from previous tests
	server.clearClients()

	// list of routes to bind the running trams to
	tests := []int{1, 96, 101, 1, 96, 109}
	// initialise goroutes for client calls
	ch := make(chan *Tram)
	for i, _ := range tests {

		go func(route int) {
			worker := new(Client)
			worker.Init("localhost:1234")
			Logger(fmt.Sprintf("tram %s initiated", worker.TramObj.TramID.String()))
			err := worker.RegisterRoute(route)
			if err != nil {
				fmt.Println("oh oh")
			}
			Logger(fmt.Sprintf("tram %s registered to route %d", worker.TramObj.TramID.String(), route))
			// run tram through 10 stops
			for j := 0; j < 10; j++ {
				worker.AsyncAdvance()
				ch <- worker.TramObj
			}

		}(tests[i])
	}
	// listen for messages from clients, until a 30sec gap of inactivity
	for {
		select {
		case response := <-ch:
			fmt.Printf("Tram %s is at tram stop: %d\n", response.TramID.String(), response.CurrentStop)
			Logger(fmt.Sprintf("tram %s arrived at stop %d", response.TramID.String(), response.CurrentStop))
		case <-time.After(30 * 1e9):
			// 30 sec timeout
			Logger(fmt.Sprintf("loop timeout, no events for 30 seconds"))
			return
		}
	}
}
