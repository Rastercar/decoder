package tcp

import (
	"net"
)

var MAX_INVALID_MESSAGES_PER_CONN = 10

func handlePackets(packets []byte) ([]byte, error) {
	return nil, nil
}

// HandleRequest deals with the connection between tracker and decoder,
// listening to tracker packets until the connection is dropped or the
// too many invalid packets are recieved.
func HandleRequest(c net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)

		if err != nil {
			c.Close()
			break
		}

		req := buf[:n]

		res, err := handlePackets(req)

		if err != nil {
			c.Close()
			break
		}

		if res != nil {
			c.Write(res)
		}
	}
}
