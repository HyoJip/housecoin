package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
	"strings"
)

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
}

const difficultyLevel = 2

var ErrNotFound = errors.New("Block Not Found")

func CreateBlock(data string, prevHash string, height int) *Block {
	block := &Block{data, "", prevHash, height, difficultyLevel, 0}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.mine()
	block.persist()
	return block
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		blockAsString := fmt.Sprint(b)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))
		fmt.Printf("Block As String:%s\nHash:%s\nNonce:%d\n\n\n", blockAsString, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}

	}
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
