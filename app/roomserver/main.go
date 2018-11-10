package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"hhserver/app/roomserver/client"
	"hhserver/app/roomserver/localcommand"
	"net/http"
	"strconv"
)

func main() {
	go localcommand.Local_manager_start()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)
	wsAddr := ":" + strconv.Itoa(8491)

	fmt.Println("Server Started")

	err := http.ListenAndServe(wsAddr, mux)



	if err != nil {
		fmt.Println("err : ",err)
		//panic(err)
	}
}






func wsHandler(res http.ResponseWriter, req *http.Request){
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	c := client.NewClient(conn)
	c.ServeClient()
}






