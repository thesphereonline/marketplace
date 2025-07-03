package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prev_hash"`
	Hash         string        `json:"hash"`
	Validator    string        `json:"validator"`
}

func (b *Block) CalculateHash() string {
	data, _ := json.Marshal(b.Transactions)
	raw := string(b.Index) + string(b.Timestamp) + string(data) + b.PrevHash + b.Validator
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
