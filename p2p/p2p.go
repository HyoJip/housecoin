package p2p

import (
	"fmt"
	"github.com/Hyojip/housecoin/utils"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(writer http.ResponseWriter, request *http.Request) {
	address := utils.Splitter(request.RemoteAddr, ":", 0)
	openPort := request.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return address != "" && openPort != ""
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	utils.HandleError(err)

	initPeer(conn, address, openPort)
}

func AddPeer(address, port, openPort string) {
	fullUrl := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:])
	conn, _, err := websocket.DefaultDialer.Dial(fullUrl, nil)
	utils.HandleError(err)

	initPeer(conn, address, port)
}
