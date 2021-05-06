package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	conn  *websocket.Conn
	write chan interface{}
}

func NewWebsocketClient(conn *websocket.Conn) WebsocketClient {
	return WebsocketClient{
		conn,
		make(chan interface{}),
	}
}

func (ws *WebsocketClient) writer(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		case data := <-ws.write:
			if err := ws.conn.WriteJSON(data); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (ws *WebsocketClient) reader() {
	for {
		if _, _, err := ws.conn.ReadMessage(); err != nil {
			log.Println(err)
			return
		}
	}
}

func (ws *WebsocketClient) Handler() {
	quit := make(chan bool)
	go ws.writer(quit)
	ws.reader()
	quit <- true
}
