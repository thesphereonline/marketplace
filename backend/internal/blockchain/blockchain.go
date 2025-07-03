package blockchain

import (
	"strings"
	"sync"
	"time"

	"github.com/thesphereonline/marketplace/internal/core"
)

type Blockchain struct {
	Chain    []core.Block
	Mempool  []core.Transaction
	Listings map[string]core.Transaction // NFT ID → listing TX
	lock     sync.Mutex
}

func NewBlockchain() *Blockchain {
	genesis := core.Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []core.Transaction{},
		PrevHash:     "",
		Hash:         "",
		Validator:    "genesis",
	}
	genesis.Hash = genesis.CalculateHash()
	return &Blockchain{
		Chain:    []core.Block{genesis},
		Mempool:  []core.Transaction{},
		Listings: make(map[string]core.Transaction),
	}
}

func (bc *Blockchain) AddTransaction(tx core.Transaction) {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.Mempool = append(bc.Mempool, tx)
}

func (bc *Blockchain) GetLatestBlock() core.Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) MineBlock(validator string) core.Block {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	prev := bc.GetLatestBlock()
	newBlock := core.Block{
		Index:        prev.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: bc.Mempool,
		PrevHash:     prev.Hash,
		Validator:    validator,
	}
	newBlock.Hash = newBlock.CalculateHash()
	bc.Chain = append(bc.Chain, newBlock)

	// Process new listings
	for _, tx := range bc.Mempool {
		if strings.HasPrefix(tx.Data, "LIST:") {
			bc.Listings[tx.ID] = tx
		}
	}

	bc.Mempool = []core.Transaction{}
	return newBlock
}
