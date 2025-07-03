package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/thesphereonline/marketplace/internal/blockchain"
	"github.com/thesphereonline/marketplace/internal/core"
)

func parseArgs() (string, string) {
	var owner string
	var meta string

	for i, arg := range os.Args {
		if arg == "--owner" && i+1 < len(os.Args) {
			owner = os.Args[i+1]
		}
		if arg == "--meta" && i+1 < len(os.Args) {
			meta = os.Args[i+1]
		}
	}
	return owner, meta
}

func main() {
	if len(os.Args) < 5 || os.Args[1] != "--mint" {
		fmt.Println("❌ Usage: go run ./cmd/nft/main.go --mint --owner <ADDR> --meta <META>")
		return
	}

	owner, meta := parseArgs()

	if owner == "" || meta == "" {
		fmt.Println("❌ Missing required --owner or --meta value.")
		return
	}

	tx := core.Transaction{
		ID:     "tx-nft-" + strings.ReplaceAll(meta, " ", "_"),
		From:   "system",
		To:     owner,
		Amount: 0,
		Data:   meta,
	}

	chain := blockchain.NewBlockchain()
	chain.AddTransaction(tx)

	block := chain.MineBlock("nft-minter")
	fmt.Println("✅ NFT Minted!")
	fmt.Println("📦 Block Hash:", block.Hash)
	fmt.Println("🎨 Metadata:", meta)
}
