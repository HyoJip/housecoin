package p2p

import (
	"encoding/json"
	"fmt"
	"github.com/Hyojip/housecoin/blockchain"
	"github.com/Hyojip/housecoin/utils"
)

type MessageKind int

type Message struct {
	Kind    MessageKind
	Payload []byte
}

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlockRequest
	MessageAllBlockResponse
	MessageNotifyNewBlock
	MessageNotifyNewTx
)

func sendNewestBlock(p *peer) {
	block, err := blockchain.FindBlock(blockchain.GetBlockchain().NewestHash)
	utils.HandleError(err)
	m := makeMessage(MessageNewestBlock, block)
	p.inbox <- m
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	message := utils.ToJSON(m)
	return message
}
func handleMessage(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received newest block from %s\n", p.key)
		var newBlock blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &newBlock))
		serverNewestBlock, err := blockchain.FindBlock(blockchain.GetBlockchain().NewestHash)
		utils.HandleError(err)
		if serverNewestBlock.Height <= newBlock.Height {
			fmt.Printf("Requesting all blocks from %s\n", p.port)
			requestAllBlock(p)
		} else {
			fmt.Printf("Sending Newest blocks from %s\n", p.port)
			sendNewestBlock(p)
		}
		fmt.Printf("Peer: %s, Block: %+v\n", p.port, newBlock)
	case MessageAllBlockRequest:
		fmt.Printf("%s wants all blocks\n", p.port)
		sendAllBlockMessage(p)
	case MessageAllBlockResponse:
		fmt.Printf("Received All blocks From %s\n", p.port)
		var blocks []*blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &blocks))
		blockchain.GetBlockchain().Replace(blocks)
	case MessageNotifyNewBlock:
		fmt.Printf("Notified block From %s\n", p.port)
		var block *blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &block))
		blockchain.GetBlockchain().AddPeerBlock(block)
	case MessageNotifyNewTx:
		fmt.Printf("Notified Transaction From %s\n", p.port)
		var tx *blockchain.Tx
		utils.HandleError(json.Unmarshal(m.Payload, &tx))
		blockchain.GetMempool().AddPeerTx(tx)
	}
}

func requestAllBlock(p *peer) {
	m := makeMessage(MessageAllBlockRequest, nil)
	p.inbox <- m
}

func sendAllBlockMessage(p *peer) {
	m := makeMessage(MessageAllBlockResponse, blockchain.FindBlocks())
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNotifyNewBlock, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNotifyNewTx, tx)
	p.inbox <- m
}
