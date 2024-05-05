package customwebsocket

import "fmt"

type Pool struct {
	Register chan *Client

	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("total connection pool:- ", len(pool.Clients))
			for key, _ := range pool.Clients {
				fmt.Println(key)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Join"})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("total connection pool:- ", len(pool.Clients))
			for key, _ := range pool.Clients {
				fmt.Println(key)
				client.Conn.WriteJSON(Message{Type: 1, Body: "User disconnected"})
			}
		case message := <-pool.Broadcast:
			fmt.Println("bradcasting a message")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
