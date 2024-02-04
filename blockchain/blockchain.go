package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

type blockchain struct {
	blocks []*Block
}

var chain *blockchain
var once sync.Once
var ErrNotFound error = errors.New("block Not Found")

func GetBlockchain() *blockchain {
	if chain == nil {
		once.Do(func() {
			chain = &blockchain{}
			chain.AddBlock("Genesis Block")
		})
	}
	return chain
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(AllBlock()) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func getLastHash() string {
	chain := GetBlockchain()
	blocks := chain.blocks
	if l := len(blocks); l > 0 {
		return blocks[l-1].Hash
	}
	return ""
}

func AllBlock() []*Block {
	return GetBlockchain().blocks
}

func GetBlock(height int) (*Block, error) {
	if length := len(chain.blocks); height < 1 || length < height {
		return nil, ErrNotFound
	}
	return chain.blocks[height-1], nil
}
