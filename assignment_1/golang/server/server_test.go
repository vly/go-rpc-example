package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test whether server successfully binds to a port
func TestPortbind(t *testing.T) {
	obj, ok := ServerBind()
	defer obj.Close()

	assert.Nil(t, ok, "Error on bind")
	assert.NotNil(t, obj, "Server object empty")
}

func TestReceiveMessage(t *testing.T) {
	obj, _ := ServerBind()
	defer obj.Close()
	ServerListen(obj)

}
