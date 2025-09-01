package main

import (
	"context"
	"flag"
	"log"

	"example.com/remote/client/agent"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "control plane address")
	token := flag.String("token", "", "enrollment token")
	flag.Parse()

	a := agent.New()
	if _, err := a.Run(context.Background(), *addr, *token); err != nil {
		log.Fatalf("agent exited: %v", err)
	}
}
