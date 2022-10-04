package tcp

import (
	"net"
	"reciever-ms/protocols/h02/decoder"
)

var dec = decoder.New(true)

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

		_, err = dec.Decode(req)

		if err != nil {
			c.Close()
			break
		}

		// TODO: FIXME
		// if res != nil {
		// 	c.Write(res)
		// }
	}
}
