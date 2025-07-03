package main

import (
	"fmt"
	"os"

	"github.com/thesphereonline/marketplace/internal/api"
	"github.com/thesphereonline/marketplace/internal/blockchain"
)

func main() {
	fmt.Println("🚀 Sphere Market full node starting...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	chain := blockchain.NewBlockchain()
	server := api.NewServer(chain)

	go server.Start(":" + port)

	select {} // keep running
}
