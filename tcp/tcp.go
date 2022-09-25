package tcp

import (
	"fmt"
	"net"
)

type ReqHandler func(net.Conn)

func Listen(address string, handler ReqHandler) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("Listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err == nil {
			go handler(conn)
		}
	}
}
