package p2p

import (
	"fmt"
	"github.com/Hyojip/housecoin/utils"
	"github.com/gorilla/websocket"
	"sync"
)

var Peers = peers{
	v: make(map[string]*peer),
	m: sync.Mutex{},
}

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

// goroutine이 총 3개??
// 메인스레드, write, read
func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		key:     key,
		address: address,
		port:    port,
		conn:    conn,
		inbox:   make(chan []byte),
	}

	go p.read()
	go p.write()

	Peers.v[key] = p
	return p
}

func (p *peer) read() {
	defer p.close()

	for {
		message := &Message{}
		err := p.conn.ReadJSON(message)
		if err != nil {
			return
		}

		handleMessage(message, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			return
		}
		utils.HandleError(p.conn.WriteMessage(websocket.TextMessage, m))
	}
}

func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	utils.HandleError(p.conn.Close())
	delete(Peers.v, p.key)
}

func FindPeers(p *peers) []string {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	var keys []string
	for k := range p.v {
		keys = append(keys, k)
	}
	return keys
}
