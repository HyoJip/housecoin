package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

var ErrNotFound = errors.New("Block Not Found")

func CreateBlock(data string, prevHash string, height int) *Block {
	block := &Block{data, "", prevHash, height}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockB := db.Block(hash)
	if blockB == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockB)
	return block, nil
}
