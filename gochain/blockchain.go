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

type BlockChain struct {
	transactionPool []string
	chain 			[]*Block
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash string) *Block {
	new_block := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, new_block)

	return new_block
}

func NewBlockChain() *BlockChain {
	bc := new(BlockChain) 
	bc.CreateBlock(0, "init hash")

	return bc
}

func (bc *BlockChain) Print()  {
	for i, block := range bc.chain {
		fmt.Printf("Chain %d\n", i)
		block.Print()
	}
}

func init() {
	log.SetPrefix("GoChain: ")
}

func main() {
	blockChain := NewBlockChain()
	blockChain.Print()
}
