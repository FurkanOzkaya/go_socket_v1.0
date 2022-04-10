package websocket

import (
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
}

type Pool struct {
	Operation    chan *ConnectionModel
	Clients map[string]*websocket.Conn
}

func NewPool() *Pool {
	return &Pool{
		Operation:    make(chan *ConnectionModel),
		Clients:  make(map[string]*websocket.Conn),
	}
}

func (pool *Pool) Start() {

	for {
		client := <-pool.Operation
		switch client.Operation{
		case CONNECT:
			if pool.Clients == nil {
			    pool.Clients = make(map[string]*websocket.Conn)
			}
			pool.Clients[client.User] = client.Conn
			send_msg := ConnectionModel{Operation: CONNECT, Message: "success"}
			go sendMessage(client.Conn, send_msg)
		case MESSAGE:
			log.Println("[Start] pool.MESSAGE chan to message sendMessage")
			log.Println(client)
			targetUser := pool.Clients[client.To]
			//sendMsg := ConnectionModel{Operation: MESSAGE, From: client.From, To: client.To, Message: client.Message}
			go sendMessage(targetUser,  *client)
		case DISCONNECT:
			log.Println(client)
			delete(pool.Clients, client.User)
			send_msg :=ConnectionModel{Operation: DISCONNECT, Message: "success"}
			go sendMessage(client.Conn,  send_msg)

		}
	}
}

func sendMessage(conn *websocket.Conn, msg ConnectionModel){
    if err := conn.WriteJSON(msg); err != nil {
        log.Println(err)
        return
    }
}

