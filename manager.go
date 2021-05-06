package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Manager struct {
	scale      *Scale
	upgrader   websocket.Upgrader
	websockets map[*WebsocketClient]bool
	weight     float64
	connected  bool
}

func NewManager(scale *Scale) *Manager {
	return &Manager{
		scale,
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		map[*WebsocketClient]bool{},
		0,
		false,
	}
}

func (m *Manager) broadcast(data interface{}) {
	for client := range m.websockets {
		select {
		case client.write <- data:
		default:
			delete(m.websockets, client)
		}
	}
}

func (m *Manager) Run() {
	for {
		select {
		case weight := <-m.scale.weight:
			m.weight = weight
			m.broadcast(map[string]interface{}{
				"type": "weight",
				"data": weight,
			})
		case connected := <-m.scale.connected:
			m.connected = connected
			m.broadcast(map[string]interface{}{
				"type": "connected",
				"data": connected,
			})
		}
	}
}

func (m *Manager) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := m.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer ws.Close()

		log.Println("Client connected")
		client := NewWebsocketClient(ws)
		m.websockets[&client] = true

		go func() {
			client.write <- map[string]interface{}{
				"type": "connected",
				"data": m.connected,
			}
		}()

		client.Handler()
		delete(m.websockets, &client)
	}
}
