package main

import (
	"fmt"
	"log"
	"os"
	"reciever-ms/tcp"
	"reciever-ms/tracer"
)

func initTracing() {
	cfg := tracer.TracerConfig{
		ServiceName:    "positions-ms",
		ExportEndpoint: "http://localhost:14268/api/traces",
	}

	if err := tracer.SetGlobalTracer(&cfg); err != nil {
		log.Fatalf("[TRACER] failed to init tracer: %v", err)
	}
}

func main() {
	initTracing()

	err := tcp.Listen("localhost:3003", tcp.HandleRequest)
	if err != nil {
		fmt.Printf("\n[SERVER] Failed to start: %v", err)
		os.Exit(1)
	}
}
