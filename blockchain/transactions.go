package blockchain

import (
	"errors"
	"fmt"
	"github.com/Hyojip/housecoin/utils"
	"github.com/Hyojip/housecoin/wallet"
	"sync"
	"time"
)

const (
	minerReward int = 50
)

var ErrorNotEnoughMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("invalid transaction")

var m *Mempool
var memOnce sync.Once

type Mempool struct {
	Txs map[string]*Tx
	M   sync.Mutex
}

func GetMempool() *Mempool {
	if m == nil {
		memOnce.Do(func() {
			m = &Mempool{
				Txs: make(map[string]*Tx),
				M:   sync.Mutex{},
			}
		})
	}
	return m
}

func (m *Mempool) AddPeerTx(tx *Tx) {
	m.M.Lock()
	defer m.M.Unlock()

	m.Txs[tx.Id] = tx
}

type Tx struct {
	Id        string   `json:"id,omitempty"`
	Timestamp int      `json:"timestamp,omitempty"`
	TxIns     []*TxIn  `json:"txIns,omitempty"`
	TxOuts    []*TxOut `json:"txOuts,omitempty"`
}

func (t *Tx) generateId() {
	t.Id = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, tx := range t.TxIns {
		tx.Signature = wallet.Sign(wallet.Wallet(), t.Id)
	}
}

func validate(t *Tx) bool {
	for _, txIn := range t.TxIns {
		prevTx := FindTx(txIn.TxId)
		if prevTx == nil {
			return false
		}

		address := prevTx.TxOuts[txIn.Index].Address
		if isCorrect := wallet.Verify(txIn.Signature, t.Id, address); !isCorrect {
			return false
		}
	}
	return true
}

type TxIn struct {
	TxId      string `json:"txId,omitempty"`
	Index     int    `json:"index,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type TxOut struct {
	Address string `json:"address,omitempty"`
	Amount  int    `json:"amount,omitempty"`
}

type UTxOut struct {
	TxID   string `json:"txID,omitempty"`
	Index  int    `json:"index"`
	Amount int    `json:"amount,omitempty"`
}

func (m *Mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	m.Txs[tx.Id] = tx
	return tx, nil
}

func (m *Mempool) confirmTx() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}

func containsTx(m *Mempool, out *UTxOut) bool {
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
		{"", -1, "COINBASE"},
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
		return nil, ErrorNotEnoughMoney
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
			TxId:      uOut.TxID,
			Index:     uOut.Index,
			Signature: from,
		}
		txIns = append(txIns, in)
		total += uOut.Amount
	}

	if change := total - amount; change != 0 {
		changeOut := &TxOut{
			Address: from,
			Amount:  change,
		}
		txOuts = append(txOuts, changeOut)
	}
	out := &TxOut{
		Address: to,
		Amount:  amount,
	}
	txOuts = append(txOuts, out)

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.generateId()
	tx.sign()
	if !validate(tx) {
		return nil, ErrorNotValid
	}
	return tx, nil
}
