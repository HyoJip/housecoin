package blockchain

import (
	"errors"
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
	"strings"
	"time"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"getDifficulty"`
	Nonce        int    `json:"nonce"`
	TimeStamp    int    `json:"timeStamp"`
	Transactions []*Tx  `json:"transactions"`
}

var ErrNotFound = errors.New("block not found")

func CreateBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: getDifficulty(chain),
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.confirmTx()
	persistBlock(block)
	return block
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.TimeStamp = int(time.Now().Unix())
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}

	}
}

func persistBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func restoreBlock(b *Block, data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockB := db.Block(hash)
	if blockB == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	restoreBlock(block, blockB)
	return block, nil
}
