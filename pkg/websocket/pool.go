package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)



type ConnectionModel struct {
	Operation string `json:"operation"`
	User      string `json:"user"`
	From      string `json:"from"`
	To        string `json:"to"`
	Message   string `json:"message"`
	Conn	 *websocket.Conn
	MessageType int `json:"message_type"`
}

type Pool struct {
	CONNECT    chan *ConnectionModel
	MESSAGE    chan *ConnectionModel
	DISCONNECT chan *ConnectionModel
	Clients map[string]*websocket.Conn
}

func NewPool() *Pool {
	return &Pool{
		CONNECT:    make(chan *ConnectionModel),
		MESSAGE:    make(chan *ConnectionModel),
		DISCONNECT: make(chan *ConnectionModel),
		Clients:  make(map[string]*websocket.Conn),
	}
}

func (pool *Pool) Start() {

	for {
		select {
		case client := <-pool.CONNECT:
			log.Println(client, " is appending to Map of users")
			if pool.Clients == nil {
			    pool.Clients = make(map[string]*websocket.Conn)
			}
			pool.Clients[client.User] = client.Conn
			send_msg := ConnectionModel{Operation: CONNECT, Message: "success"}
			go sendMessage(client.Conn, client.MessageType, send_msg)
		case client := <-pool.MESSAGE:
			log.Println("[Start] pool.MESSAGE chan to message sendMessage")
			log.Println(client)
			targetUser := pool.Clients[client.To]
			//sendMsg := ConnectionModel{Operation: MESSAGE, From: client.From, Message: client.Message}
			go sendMessage(targetUser, client.MessageType, *client)
		case client := <-pool.DISCONNECT:
			log.Println(client)
			delete(pool.Clients, client.User)
			send_msg :=ConnectionModel{Operation: DISCONNECT, Message: "success"}
			go sendMessage(client.Conn, client.MessageType, send_msg)

		}
	}
}

func sendMessage(conn *websocket.Conn,messageType int, msg ConnectionModel){
	newMsg, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err)
    }
    if err := conn.WriteMessage(messageType, []byte(newMsg)); err != nil {
        log.Println(err)
        return
    }
}

