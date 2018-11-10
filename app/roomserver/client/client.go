package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"hhserver/app/roomserver/uid"

)

func init(){
	go Manager.start()
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

var Manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}
type ClientManager struct {
	Clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewClient(conn *websocket.Conn) *Client{
	client := &Client{id: uid.GeneratorID(), socket: conn, send: make(chan []byte)}
	Manager.register <- client
	return client
}

func (c *Client)ServeClient(){
	go c.read()
	go c.write()
}

func (c *Client)GetID() string{
	return c.id
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for c := range manager.Clients {
		if c != ignore {
			c.send <- message
		}
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case c := <-manager.register:
			manager.Clients[c] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			manager.send(jsonMessage, c)
			fmt.Printf("A new client connected : uid %s , addr : %s\n",c.id,c.socket.RemoteAddr())
		case conn := <-manager.unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.Clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}


func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) read() {
	defer func() {
		Manager.unregister <- c
		c.socket.Close()
	}()

	for {
		//messagetype 0 明文 1 binary 2 close
		messagetype , message, err := c.socket.ReadMessage()
		if err != nil {
			Manager.unregister <- c
			c.socket.Close()
			break
		}
		switch messagetype {
		case websocket.TextMessage:
			fmt.Printf("SERVER received text from %s : %s\n",c.id,string(message))
		case websocket.BinaryMessage:
			fmt.Println("SERVER received binary from %s : %s\n",c.id,string(message))
		case websocket.CloseMessage:
			fmt.Println("SERVER received close from %s : %s\n",c.id,string(message))
		}

		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		Manager.broadcast <- jsonMessage
	}
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}