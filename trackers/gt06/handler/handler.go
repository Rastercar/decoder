package handler

import (
	"net"
)

// HandleRequest deals with the connection between tracker and decoder,
// listening to tracker packets until the connection is dropped or the
// too many invalid packets are recieved.
func HandleRequest(c net.Conn) {
	s := Session{InvalidMsgCnt: 0}

	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)

		if err != nil {
			c.Close()
			break
		}

		req := buf[:n]

		res, err := s.handlePackets(req)

		if err != nil {
			c.Close()
			break
		}

		if res != nil {
			c.Write(res)
		}
	}
}