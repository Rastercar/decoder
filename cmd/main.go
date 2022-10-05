package main

import (
	"context"
	"log"
	"reciever-ms/config"
	h02 "reciever-ms/protocol/h02/handler"
	"reciever-ms/server"
	"reciever-ms/tracer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("[CONFIG] failed to parse config: %v", err)
	}

	err = tracer.Start(&cfg.Tracer)
	if err != nil {
		log.Fatalf("[TRACER] failed to init tracer: %v", err)
	}

	defer tracer.Stop(ctx)

	h02Handler := h02.New(cfg)

	err = server.Listen(":3003", h02Handler.HandleRequest)
	if err != nil {
		log.Fatalf("[SERVER] Failed to start: %v", err)
	}
}
