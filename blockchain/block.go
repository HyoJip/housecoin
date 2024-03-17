package blockchain

import (
	"errors"
	"fmt"
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
	"strings"
	"time"
)

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	TimeStamp  int    `json:"timeStamp"`
}

var ErrNotFound = errors.New("Block Not Found")

func CreateBlock(data string, prevHash string, height int) *Block {
	block := &Block{Data: data, Hash: "", PrevHash: prevHash, Height: height, Difficulty: chain.difficulty(), Nonce: 0}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = utils.Hash(payload)
	block.mine()
	block.persist()
	return block
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		fmt.Printf("Block As Target:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.TimeStamp = int(time.Now().Unix())
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
