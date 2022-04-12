package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
		//delete user if its broken
	}

	return conn, nil
}

const (
	CONNECT    string = "connect"
	MESSAGE    string = "message"
	DISCONNECT string = "disconnect"
)

func Reader(pool *Pool, ws *websocket.Conn) {
	for {
		var wsmessage ConnectionModel
		err := ws.ReadJSON(&wsmessage)
		if err != nil {
			log.Println(err)
			return
		}
		if wsmessage.Operation == CONNECT || wsmessage.Operation == MESSAGE || wsmessage.Operation == DISCONNECT {
			wsmessage.Conn = ws
			pool.Operation <- &wsmessage
		} else {
			log.Println("Wrong Operation Type")
		}

		fmt.Printf("Message Received: %+v\n", wsmessage)

	}
}
