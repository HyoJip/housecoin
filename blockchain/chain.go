package blockchain

import (
	"encoding/json"
	"github.com/Hyojip/housecoin/db"
	"github.com/Hyojip/housecoin/utils"
	"net/http"
	"sync"
)

const (
	defaultDifficulty   int = 2
	difficultyInterval  int = 5
	blockInterval       int = 2
	allowedRangeMinutes int = 2
)

var once sync.Once
var chain *blockchain

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

func (b *blockchain) AddBlock() *Block {
	block := CreateBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
	return block
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	newestBlock := newBlocks[0]
	b.Height = len(newBlocks)
	b.NewestHash = newestBlock.Hash
	b.CurrentDifficulty = newestBlock.Difficulty
	persistBlockchain(b)
	db.EmptyBlock()

	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.M.Lock()
	defer b.m.Unlock()
	defer m.M.Unlock()

	b.Height += 1
	b.NewestHash = newBlock.Hash
	b.CurrentDifficulty = newBlock.Difficulty

	persistBlockchain(b)
	persistBlock(newBlock)

	for _, t := range newBlock.Transactions {
		if _, ok := m.Txs[t.Id]; ok {
			delete(m.Txs, t.Id)
		}
	}
}

func GetBlockchain() *blockchain {
	once.Do(func() {
		chain = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			chain.AddBlock()
		} else {
			restoreBlockchain(chain, checkpoint)
		}
	})

	return chain
}

func FindBlocks() []*Block {
	b := GetBlockchain()
	b.m.Lock()
	defer b.m.Unlock()

	var blocks []*Block
	hashCursor := b.NewestHash
	for hashCursor != "" {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)

		hashCursor = block.PrevHash
	}
	return blocks
}

func FindUTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, b := range FindBlocks() {
		for _, t := range b.Transactions {
			for _, in := range t.TxIns {
				if in.Signature == "COINBASE" {
					break
				}
				if FindTx(in.TxId).TxOuts[in.Index].Address == address {
					creatorTxs[in.TxId] = true
				}
			}
			for idx, out := range t.TxOuts {
				if out.Address == address {
					if _, ok := creatorTxs[t.Id]; !ok {
						uTxOut := &UTxOut{
							TxID:   t.Id,
							Index:  idx,
							Amount: out.Amount,
						}
						if !containsTx(GetMempool(), uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func FindBalanceByAddress(address string) int {
	total := 0
	for _, txOut := range FindUTxOutsByAddress(address) {
		total += txOut.Amount
	}
	return total
}

func FindTxs() []*Tx {
	var txs []*Tx
	for _, tx := range FindBlocks() {
		txs = append(txs, tx.Transactions...)
	}
	return txs
}

func FindTx(txId string) *Tx {
	for _, tx := range FindTxs() {
		if txId == tx.Id {
			return tx
		}
	}
	return nil
}

func restoreBlockchain(b *blockchain, data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func recalculateDifficulty(b *blockchain) int {
	blocks := FindBlocks()
	recentBlock := blocks[0]
	lastDifficultyBlock := blocks[difficultyInterval-1]
	diffMinutes := (recentBlock.TimeStamp / 60) - (lastDifficultyBlock.TimeStamp / 60)
	expectedMinutes := difficultyInterval * blockInterval
	if diffMinutes <= expectedMinutes-allowedRangeMinutes {
		b.CurrentDifficulty++
	} else if diffMinutes >= expectedMinutes+allowedRangeMinutes && b.CurrentDifficulty > 0 {
		b.CurrentDifficulty--
	}
	return b.CurrentDifficulty
}

func Status(b *blockchain, w http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	utils.HandleError(json.NewEncoder(w).Encode(b))
}
