package main

import (
	"fmt"
	"os"
	"reciever-ms/tcp"
)

func main() {
	err := tcp.Listen("localhost:3003", tcp.HandleRequest)
	if err != nil {
		fmt.Printf("\n[SERVER] Failed to start: %v", err)
		os.Exit(1)
	}
}
