package p2p

import (
	"fmt"
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

func initPeer(conn *websocket.Conn, address, port string) *peer {
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
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("%s\n", message)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			return
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
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
