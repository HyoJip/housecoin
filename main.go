package main

import (
	"fmt"
	"github.com/Hyojip/housecoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()

	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")

	for _, b := range blockchain.AllBlock() {
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %s\n", b.Hash)
		fmt.Printf("PrevHash: %s\n", b.PrevHash)
	}
}
