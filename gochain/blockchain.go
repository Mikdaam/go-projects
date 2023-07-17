package main

import "log"

type Block struct {
	nonce int
	previousHash string
	timestamp int64
	transactions []string
}

func init() {
	log.SetPrefix("GoChain: ")
}

func main() {
	log.Println("test")
}
