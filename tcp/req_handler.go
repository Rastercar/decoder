package tcp

import (
	"fmt"
	"net"
	"reciever-ms/trackers/gt06"
)

var decoder = gt06.NewDecoder(true)

func HandleRequest(c net.Conn) {
	for {
		buf := make([]byte, 1024)

		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("\nError reading:", err.Error())
			c.Close()
			return
		}

		req := buf[:n]

		res, err := decoder.Decode(req)
		if err != nil {
			fmt.Println("Error decoding:", err.Error())
		}

		c.Write(res)
	}
}
