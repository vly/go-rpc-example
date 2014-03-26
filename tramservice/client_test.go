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

	// clear current clients from previous tests
	server.clearClients()

	// list of routes to bind the running trams to
	tests := []int{1, 96, 1}

	// initialise test trams, ready to start moving
	// workingClients := make([]Client, len(tests))
	// for i, b := range workingClients {
	// 	b.Init("localhost:1234")
	// 	err := b.RegisterRoute(tests[i])
	// 	assert.Nil(t, err, "Error registering route")
	// }
	// testChannels := make([]chan int, len(tests))
	// for i, b := range testChannels {
	// 	go func(route int) {
	// 		worker := new(Client)
	// 		worker.Init("localhost:1234")
	// 		err := worker.RegisterRoute(route)
	// 		if err != nil {
	// 			fmt.Println("oh oh")
	// 		}
	// 		worker.AsyncAdvance()
	// 		b <- 1
	// 	}(tests[i])
	// }
	// fmt.Println(len(testChannels))

	// }
	var chans [3]chan string
	for i := range chans {
		chans[i] = make(chan string)
	}
	for i, _ := range chans {
		go func(route int) {
			worker := new(Client)
			worker.Init("localhost:1234")
			err := worker.RegisterRoute(route)
			if err != nil {
				fmt.Println("oh oh")
			}
			for j := 0; j < 5; j++ {
				worker.AsyncAdvance()
				chans[i] <- fmt.Sprintf("%s: %d", worker.TramObj.TramID.String(), worker.TramObj.CurrentStop)
			}

		}(tests[i])
	}

	for {
		select {
		case <-chans[0]:
			fmt.Printf("Tram 1 received")
		case <-chans[1]:
			fmt.Println("Tram 2 received")
		case response := <-chans[2]:
			fmt.Printf("Tram 3 received: %s\n", response)
		case <-time.After(5 * 1e9):
			// 30 sec timeout
			return
		}
	}
}
