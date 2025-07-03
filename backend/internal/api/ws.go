package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/thesphereonline/marketplace/internal/blockchain"
)

type WebSocketHub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
	chain   *blockchain.Blockchain
}

func NewWebSocketHub(chain *blockchain.Blockchain) *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[*websocket.Conn]bool),
		chain:   chain,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins
	},
}

func (hub *WebSocketHub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("❌ WS upgrade failed:", err)
		return
	}

	hub.lock.Lock()
	hub.clients[conn] = true
	hub.lock.Unlock()

	go hub.listen(conn)
}

func (hub *WebSocketHub) listen(conn *websocket.Conn) {
	defer func() {
		hub.lock.Lock()
		delete(hub.clients, conn)
		hub.lock.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (hub *WebSocketHub) Broadcast(message string) {
	hub.lock.Lock()
	defer hub.lock.Unlock()

	for client := range hub.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close()
			delete(hub.clients, client)
		}
	}
}
