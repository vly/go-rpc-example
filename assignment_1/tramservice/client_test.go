package tramservice

import (
	//	"fmt"
	//"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// init server if hasn't been already
var server Server

func ConnectServer() {
	if server.Status != true {
		go server.Init()
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
		b.RegisterRoute(tests[i][0])
		b.SetCurrentLocation(tests[i][1], tests[i][2])
		nextStop, err := b.GetNextStop()
		assert.Nil(t, err, "Error getting next stop.")
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

// func TestSequentialPathing(t *testing.T) {
// 	// init server
// 	ConnectServer()

// 	// set of trams w
// 	tramID, _ := uuid.NewV4()
// 	tests := []Tram{
// 		Tram{tramID, 1, 0, 2},
// 		Tram{tramID, 96, 0, 1},
// 		Tram{tramID, 109, 0, 2},
// 		Tram{tramID, 101, 0, 4},
// 	}
// 	workingClients := make([]*Client, len(tests))
// 	for a, b := range tests {
// 		workingClients[a] = new(Client)
// 		err := workingClients[a].Init("localhost:1234")
// 		assert.Nil(t, err, "Error initialising client")
// 		result, err := workingClients[a].UpdateTramLocation(&b)
// 		assert.Nil(t, err, "Error updating tram location")
// 		assert.NotNil(t, result, "Error updating tram location")
// 		workingClients[a].AdvanceTram(&b)

// 	}
// }