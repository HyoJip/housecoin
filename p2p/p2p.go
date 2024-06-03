package p2p

import (
	"fmt"
	"github.com/Hyojip/housecoin/blockchain"
	"github.com/Hyojip/housecoin/utils"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(writer http.ResponseWriter, request *http.Request) {
	address := utils.Splitter(request.RemoteAddr, ":", 0)
	openPort := request.URL.Query().Get("openPort")
	fmt.Printf("%s wants to upgrade\n", openPort)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return address != "" && openPort != ""
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	utils.HandleError(err)

	initPeer(conn, address, openPort)
}

func AddPeer(address, port, openPort string) {
	fmt.Printf("%s wants to connect to port %s\n", openPort[1:], port)
	fullUrl := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:])
	conn, _, err := websocket.DefaultDialer.Dial(fullUrl, nil)
	utils.HandleError(err)

	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}
