package blockchain

import (
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var chain *blockchain
var once sync.Once

func GetBlockchain() *blockchain {
	if chain == nil {
		once.Do(func() {
			chain = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				chain.AddBlock("Genesis Block")
			} else {
				chain.restore(checkpoint)
			}
		})
	}
	return chain
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock(data string) {
	block := CreateBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func FindBlocks() []*Block {
	var blocks []*Block
	hashCursor := GetBlockchain().NewestHash
	for hashCursor != "" {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)

		hashCursor = block.PrevHash
	}
	return blocks
}
