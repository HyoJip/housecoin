package blockchain

import (
	"errors"
	"github.com/Hyojip/housecoin/utils"
	"time"
)

const (
	minerReward int = 50
)

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id,omitempty"`
	Timestamp int      `json:"timestamp,omitempty"`
	TxIns     []*TxIn  `json:"txIns,omitempty"`
	TxOuts    []*TxOut `json:"txOuts,omitempty"`
}

func (t *Tx) generateId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	TxId  string `json:"txId,omitempty"`
	Index int    `json:"index,omitempty"`
	Owner string `json:"owner,omitempty"`
}

type TxOut struct {
	Owner  string `json:"owner,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

type UTxOut struct {
	TxID   string `json:"txID,omitempty"`
	Index  int    `json:"index"`
	Amount int    `json:"amount,omitempty"`
}

type mempool struct {
	Txs []*Tx
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("house", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) confirmTx() []*Tx {
	coinbase := makeCoinbaseTx("house")
	txs := append(m.Txs, coinbase)
	m.Txs = nil
	return txs
}

func containsTx(m *mempool, out *UTxOut) bool {
	isContains := false
Outer: // label
	for _, t := range m.Txs {
		for _, in := range t.TxIns {
			if in.TxId == out.TxID && in.Index == out.Index {
				isContains = true
				break Outer
			}
		}
	}
	return isContains
}

func makeCoinbaseTx(address string) *Tx {
	txIn := []*TxIn{
		{"", -1, "CODEBASE"},
	}
	txOut := []*TxOut{
		{address, minerReward},
	}

	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIn,
		TxOuts:    txOut,
	}

	tx.generateId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if FindBalanceByAddress(from) < amount {
		return nil, errors.New("Not enough money")
	}

	var txIns []*TxIn
	var txOuts []*TxOut
	uTxOuts := FindUTxOutsByAddress(from)
	total := 0
	for _, uOut := range uTxOuts {
		if total >= amount {
			break
		}
		in := &TxIn{
			TxId:  uOut.TxID,
			Index: uOut.Index,
			Owner: from,
		}
		txIns = append(txIns, in)
		total += uOut.Amount
	}

	if change := total - amount; change != 0 {
		changeOut := &TxOut{
			Owner:  from,
			Amount: change,
		}
		txOuts = append(txOuts, changeOut)
	}
	out := &TxOut{
		Owner:  to,
		Amount: amount,
	}
	txOuts = append(txOuts, out)

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.generateId()

	return tx, nil
}
