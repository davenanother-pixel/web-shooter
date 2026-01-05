package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type Client struct{ send chan []byte }
type Hub   struct{
	clients map[*Client]bool
	broadcast chan []byte
	register  chan *Client
	unregister chan *Client
}

var h = Hub{
	clients: make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (h *Hub) run(){
	for {
		select {
		case c := <-h.register:    h.clients[c] = true
		case c := <-h.unregister:  delete(h.clients, c); close(c.send)
		case m := <-h.broadcast:
			for c := range h.clients { select { case c.send <- m: default: close(c.send); delete(h.clients, c) } }
		}
	}
}

var up = websocket.Upgrader{CheckOrigin: func(r *http.Request)bool{return true}}

func wsHandler(w http.ResponseWriter, r *http.Request){
	conn, _ := up.Upgrade(w,r,nil)
	defer conn.Close()
	c := &Client{send: make(chan []byte, 256)}
	h.register <- c
	defer func(){ h.unregister <- c }()

	// simple physics loop for this client
	type S struct{ X,Y,R float64; RGB [3]int }
	s := S{400,300,20,[3]int{0,255,255}}
	go func(){
		for msg := range c.send { conn.WriteMessage(websocket.TextMessage, msg) }
	}()
	for {
		_, data, _ := conn.ReadMessage()
		var m map[string]string; json.Unmarshal(data, &m)
		switch m["cmd"] {
		case "U": s.Y -= 10
		case "D": s.Y += 10
		case "L": s.X -= 10
		case "R": s.X += 10
		}
		out, _ := json.Marshal(s)
		c.send <- out
	}
}

func main(){
	go h.run()
	http.HandleFunc("/ws", wsHandler)
	log.Println("gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
