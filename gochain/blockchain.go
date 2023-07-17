package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.nonce = nonce
	b.timestamp = time.Now().UnixNano()
	b.previousHash = previousHash

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:      %d\n", b.timestamp)
	fmt.Printf("nonce:          %d\n", b.nonce)
	fmt.Printf("previousHash:	%s\n", b.previousHash)
	fmt.Printf("transactions:	%s\n", b.transactions)
}

func init() {
	log.SetPrefix("GoChain: ")
}

func main() {
	block := NewBlock(0, "initial hash")
	block.Print()
}
