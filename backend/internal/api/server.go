package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/thesphereonline/marketplace/internal/blockchain"
	"github.com/thesphereonline/marketplace/internal/core"
	"github.com/thesphereonline/marketplace/internal/infra"
)

type APIServer struct {
	Chain *blockchain.Blockchain
	WS    *WebSocketHub
}

func NewServer(chain *blockchain.Blockchain) *APIServer {
	return &APIServer{
		Chain: chain,
		WS:    NewWebSocketHub(chain),
	}
}

func (s *APIServer) Start(address string) {
	infra.InitPostgres()

	http.HandleFunc("/blocks", withCORS(s.handleBlocks))
	http.HandleFunc("/tx", withCORS(s.handleTx))
	http.HandleFunc("/mint", withCORS(s.handleMint))
	http.HandleFunc("/wallet/", withCORS(s.handleWallet))
	http.HandleFunc("/list", withCORS(s.handleList))
	http.HandleFunc("/buy", withCORS(s.handleBuy))
	http.HandleFunc("/listings", withCORS(s.handleGetListings))
	http.HandleFunc("/stream", s.WS.HandleConnection) // WebSocket doesn’t need CORS

	fmt.Println("🌐 API server running at", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func (s *APIServer) handleBlocks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.Chain.Chain)
}

func (s *APIServer) handleTx(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var tx core.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Chain.AddTransaction(tx)
	s.WS.Broadcast("📨 new transaction added: " + tx.ID)
	w.Write([]byte("✅ transaction added"))
}

func (s *APIServer) handleMint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Owner string `json:"owner"`
		Meta  string `json:"meta"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Owner == "" || req.Meta == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid mint request"))
		return
	}

	tx := core.Transaction{
		ID:     "nft-" + req.Meta,
		From:   "system",
		To:     req.Owner,
		Amount: 0,
		Data:   req.Meta,
	}

	s.Chain.AddTransaction(tx)
	block := s.Chain.MineBlock("api-minter")

	s.WS.Broadcast("🖼️ new NFT minted: " + req.Meta)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

func (s *APIServer) handleWallet(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimPrefix(r.URL.Path, "/wallet/")
	if addr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wallet address required"))
		return
	}

	var balance int
	for _, block := range s.Chain.Chain {
		for _, tx := range block.Transactions {
			if tx.To == addr {
				balance += tx.Amount
			}
			if tx.From == addr {
				balance -= tx.Amount
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"address": addr,
		"balance": balance,
	})
}

func (s *APIServer) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var tx core.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Chain.AddTransaction(tx)

	_, err := infra.DB.Exec(
		`INSERT INTO listings (id, owner, metadata, price) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`,
		tx.ID, tx.From, strings.TrimPrefix(tx.Data, "LIST:"), tx.Amount,
	)
	if err != nil {
		log.Println("❌ DB insert error:", err)
	}

	s.WS.Broadcast("🛒 NFT listed: " + tx.Data)
	w.Write([]byte("✅ NFT listed for sale"))
}

func (s *APIServer) handleBuy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var tx core.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	listing, exists := s.Chain.Listings[tx.ID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Listing not found"))
		return
	}

	if tx.Amount < listing.Amount {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Insufficient funds"))
		return
	}

	s.Chain.AddTransaction(tx)
	s.WS.Broadcast("💸 NFT bought: " + tx.ID)
	w.Write([]byte("✅ Purchase complete"))
}

func (s *APIServer) handleGetListings(w http.ResponseWriter, r *http.Request) {
	rows, err := infra.DB.Query(`SELECT id, owner, metadata, price FROM listings ORDER BY created_at DESC`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB query failed"))
		return
	}
	defer rows.Close()

	type Listing struct {
		ID       string `json:"id"`
		Owner    string `json:"owner"`
		Metadata string `json:"metadata"`
		Price    int    `json:"price"`
	}

	var listings []Listing
	for rows.Next() {
		var l Listing
		if err := rows.Scan(&l.ID, &l.Owner, &l.Metadata, &l.Price); err == nil {
			listings = append(listings, l)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listings)
}
