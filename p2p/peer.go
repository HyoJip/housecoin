package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var Peers = make(map[string]*Peer)

type Peer struct {
	conn *websocket.Conn
}

func initPeer(conn *websocket.Conn, address, port string) {
	key := fmt.Sprintf("%s:%s", address, port)
	Peers[key] = &Peer{conn: conn}
}
