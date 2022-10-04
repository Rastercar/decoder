package handler

import (
	"net"
	"testing"
	"time"
)

func TestHandleRequest(t *testing.T) {
	t.Run("closes the connection on client close", func(t *testing.T) {
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		go func() {
			time.Sleep(time.Millisecond * 10)
			client.Close()
		}()

		HandleRequest(server)
	})

	t.Run("closes the connection after reaching MAX_INVALID_MESSAGES_PER_CONN", func(t *testing.T) {
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		go func() {
			time.Sleep(time.Millisecond * 10)
			for i := 0; i < MAX_INVALID_MESSAGES_PER_CONN; i++ {
				nonsense := []byte{1, 2, 3}
				client.Write(nonsense)
			}
		}()

		HandleRequest(server)
	})
}
