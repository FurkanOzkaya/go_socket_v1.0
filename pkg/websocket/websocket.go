package websocket

import (
	"encoding/json"
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
    CONNECT string   = "connect"
    MESSAGE  string       = "message"
    DISCONNECT string    = "disconnect"
)

func Reader(pool *Pool, ws *websocket.Conn) {
	for {
		MessageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var wsmessage ConnectionModel
        if err := json.Unmarshal(message, &wsmessage); err != nil {
            panic(err)
        }
		switch {
			case wsmessage.Operation == CONNECT:
				log.Println(wsmessage.User, " is appending to Map of users")
				wsmessage.Conn = ws
				wsmessage.MessageType = MessageType
				pool.CONNECT <- &wsmessage
			case wsmessage.Operation == MESSAGE:
				log.Println("[Reader] Message Send")
				log.Println(wsmessage)
				wsmessage.Conn = ws
				wsmessage.MessageType = MessageType
				pool.MESSAGE <- &wsmessage
			case wsmessage.Operation == DISCONNECT:
				wsmessage.Conn = ws
				wsmessage.MessageType = MessageType
				log.Println(wsmessage)
				pool.DISCONNECT <- &wsmessage
			default:
				log.Println("Wrong Operation Type")
		}
		fmt.Printf("Message Received: %+v\n", message)

	}
}