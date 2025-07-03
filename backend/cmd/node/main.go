package main

import (
	"fmt"
	"os"

	"github.com/thesphereonline/marketplace/internal/blockchain"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "--port" {
		fmt.Println("❌ Usage: go run ./cmd/node --port <PORT>")
		return
	}
	port := os.Args[2]
	address := ":" + port

	fmt.Println("🟢 Starting Sphere node on port", port)
	node := &blockchain.Node{}
	go node.Start(address)

	select {} // Keep running forever
}
