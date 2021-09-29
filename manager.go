package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Manager struct {
	scale      *Scale
	upgrader   websocket.Upgrader
	websockets map[*WebsocketClient]bool
	weight     Weight
	connected  bool
	last       time.Time
}

func NewManager(scale *Scale) *Manager {
	return &Manager{
		scale: scale,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		websockets: map[*WebsocketClient]bool{},
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
		case w := <-m.scale.weight:
			if m.weight.stable != w.stable || (m.weight.value != w.value && time.Since(m.last) > time.Millisecond*100) {
				m.weight = w
				m.broadcast(map[string]interface{}{
					"type":   "weight",
					"data":   w.value,
					"unit":   w.unit,
					"stable": w.stable,
				})
				m.last = time.Now()
			}
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
