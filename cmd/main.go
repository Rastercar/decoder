package main

import (
	"context"
	"log"
	"reciever-ms/config"
	h02 "reciever-ms/protocol/h02/handler"
	"reciever-ms/queue"
	"reciever-ms/server"
	"reciever-ms/tracer"
)

// The mailer version/build, this gets replaced at build time to the commit SHA
// with the use of linker flags. see the ldfflags on the makefile go build command

var version = "development"
var build = "development"

func init() {
	log.Println("[GIT] build:   ", build)
	log.Println("[GIT] version: ", version)
}

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

	queue := queue.New(cfg.Rmq)
	queue.Start()
	defer queue.Stop()

	h02Handler := h02.New(cfg, &queue)

	err = server.Listen(":3003", h02Handler.HandleRequest)
	if err != nil {
		log.Fatalf("[SERVER] Failed to start: %v", err)
	}
}
