package core

type Transaction struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
	Data   string `json:"data"` // Can be NFT metadata or system call
}
