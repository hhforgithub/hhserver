package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)


func main() {


	u := url.URL{Scheme: "ws", Host:"127.0.0.1:8491", Path: "/ws"}
	var dialer *websocket.Dialer

	conn, response, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Dial response : ",response)

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Printf("Client received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
		fmt.Println("send")
	}
}

