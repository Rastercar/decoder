package tcp

import (
	"net"
	"reciever-ms/trackers/gt06"
)

var (
	decoder                       = gt06.NewDecoder(true)
	MAX_INVALID_MESSAGES_PER_CONN = 10
)

// HandleRequest deals with the connection between tracker and decoder,
// listening to tracker packets until the connection is dropped or the
// too many invalid packets are recieved.
func HandleRequest(c net.Conn) {
	s := Session{Imei: "", InvalidMsgCnt: 0}

	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)

		if err != nil {
			c.Close()
			return
		}

		req := buf[:n]
		res, err := s.handlePackets(req)

		if err != nil {
			c.Close()
			return
		}

		if res != nil {
			c.Write(res)
		}
	}
}
