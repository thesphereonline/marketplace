package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func generateWallet() {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	privBytes := priv.D.Bytes()
	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)

	os.MkdirAll("wallets", os.ModePerm)
	privPath := "wallets/" + hex.EncodeToString(pubBytes) + ".key"
	err = os.WriteFile(privPath, privBytes, 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Wallet generated.")
	fmt.Println("🔑 Public Address:", hex.EncodeToString(pubBytes))
	fmt.Println("🔐 Saved to:", privPath)
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "--generate" {
		fmt.Println("❌ Usage: go run ./cmd/wallet --generate")
		return
	}
	generateWallet()
}
