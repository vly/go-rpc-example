package tramservice

import (
	// "fmt"
	"github.com/nu7hatch/gouuid"
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

func createTram(routeID int) *Tram {
	tram := new(Tram)
	tram.TramID, _ = uuid.NewV4()
	tram.CurrentStop = 0
	tram.PreviousStop = 0

	return &tram
}

// Test GetNextStop functionality
func TestGetNextStop(t *testing.T) {
	// init server
	ConnectServer()

	// in/out test input for function
	// consists of tramID, currentstop, previousstop
	tramID, _ := uuid.NewV4()
	tests := map[Tram]int{
		Tram{tramID, 1, 3, 4}:    2,
		Tram{tramID, 1, 3, 2}:    4,
		Tram{tramID, 1, 1, 2}:    2,
		Tram{tramID, 1, 5, 4}:    4,
		Tram{tramID, 101, 34, 4}: 5,
		Tram{tramID, 2, 1, 2}:    -1,
		Tram{tramID, 1, 99, 2}:   -1,
	}

	// initialise client
	client := new(Client)
	err := client.Init("localhost:1234")
	assert.Nil(t, err, "Error initialising client")

	// prep test data and run through accuracy tests
	for a, b := range tests {
		result, err := client.GetNextStop(&a)
		assert.Nil(t, err, "Error getting next stop")
		assert.Equal(t, b, result, "Received unexpected results")
	}

}

func TestUpdateTramLocation(t *testing.T) {
	// init server
	ConnectServer()

	// in/out input for function
	tramID, _ := uuid.NewV4()
	tests := map[Tram]int{
		Tram{tramID, 2, 1, 2}:    -1,
		Tram{tramID, 1, 2, 1}:    0,
		Tram{tramID, 1, 3, 2}:    0,
		Tram{tramID, 101, 34, 4}: 0,
	}

	// initialise client
	client := new(Client)
	err := client.Init("localhost:1234")
	assert.Nil(t, err, "Error initialising client")

	// prep test data and run through accuracy tests
	for a, b := range tests {
		result, err := client.UpdateTramLocation(&a)
		assert.Nil(t, err, "Error updating tram location")
		assert.Equal(t, b, result, "Received unexpected results")
	}
	server.getStats()

}

func TestSequentialPathing(t *testing.T) {
	// init server
	ConnectServer()

	// set of trams w
	tramID, _ := uuid.NewV4()
	tests := []Tram{
		Tram{tramID, 1, 0, 2},
		Tram{tramID, 96, 0, 1},
		Tram{tramID, 109, 0, 2},
		Tram{tramID, 101, 0, 4},
	}
	workingClients := make([]*Client, len(tests))
	for a, b := range tests {
		workingClients[a] = new(Client)
		err := workingClients[a].Init("localhost:1234")
		assert.Nil(t, err, "Error initialising client")
		result, err := workingClients[a].UpdateTramLocation(&b)
		assert.Nil(t, err, "Error updating tram location")
		assert.NotNil(t, result, "Error updating tram location")
		workingClients[a].AdvanceTram(&b)

	}
}
