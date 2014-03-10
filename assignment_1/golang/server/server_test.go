package server

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func ProtoClient(message string) error {
	time.Sleep(1)
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:4444")
	con, _ := net.DialUDP("udp", nil, serverAddr)
	defer con.Close()
	_, err := con.Write([]byte(message))
	return err
}

// Test whether server successfully binds to a port
func TestServerBind(t *testing.T) {
	server := new(TramServer)
	ok := server.Bind()
	defer server.socket.Close()

	assert.Nil(t, ok, "Error on bind")
	assert.NotNil(t, server.socket, "Server object empty")
}

// Test whether server successfully receives a message
func TestReceiveMessage(t *testing.T) {
	server := new(TramServer)
	server.Bind()
	defer server.socket.Close()
	go ProtoClient("testing")
	server.Listen()
}

// Test nextTramStop functionality
func TestNextTramStop(t *testing.T) {
	server := new(TramServer)

	// check first tram route
	a, err := server.RetrieveNextStop(1, 2)
	assert.Nil(t, err, "Next stop error")
	assert.Equal(t, 3, a, "Result incorrect")

	// check last tram route result
	a, err = server.RetrieveNextStop(112, 4)
	assert.Nil(t, err, "Next stop error")
	assert.Equal(t, 29, a, "Result incorrect")
}

// Test updateTramLocation() functionality
func TestUpdateTramLocation(t *testing.T) {
	server := new(TramServer)
	tramID := 1
	currentStop := 3
	err := server.UpdateTramLocation(tramID, currentStop)
	assert.Nil(t, err, "UpdateTramLocation error")

}

// Test server new client registration functionality
func TestServerSignin(t *testing.T) {

}

// Test server client signout functionality
func TestServerSignout(t *testing.T) {

}

// Tests server message processing
func TestServerSendMessage(t *testing.T) {

}

// Test server bad message response
func TestServerSendBadMessage(t *testing.T) {

}

// Test server check for new messages functionality
func TestServerCheckmailbox(t *testing.T) {

}

// Test server message retrieval functionality
func TestServerRetrieve(t *testing.T) {

}

// Test server lookup of current registered clients functionality
func TestServerWho(t *testing.T) {

}
