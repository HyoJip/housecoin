package blockchain

import (
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
	"sync"
)

const (
	defaultDifficulty   int = 2
	difficultyInterval  int = 5
	blockInterval       int = 2
	allowedRangeMinutes int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var chain *blockchain
var once sync.Once

func GetBlockchain() *blockchain {
	if chain == nil {
		once.Do(func() {
			chain = &blockchain{
				Height: 0,
			}
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
	b.CurrentDifficulty = block.Difficulty
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

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) recalculateDifficulty() int {
	blocks := FindBlocks()
	recentBlock := blocks[0]
	lastDifficultyBlock := blocks[difficultyInterval-1]
	diffMinutes := (recentBlock.TimeStamp / 60) - (lastDifficultyBlock.TimeStamp / 60)
	expectedMinutes := difficultyInterval * blockInterval
	if diffMinutes <= expectedMinutes-allowedRangeMinutes {
		b.CurrentDifficulty++
	} else if diffMinutes >= expectedMinutes+allowedRangeMinutes {
		b.CurrentDifficulty--
	}
	return b.CurrentDifficulty
}
