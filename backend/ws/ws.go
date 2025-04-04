package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // อนุญาตทุก origin (สำหรับ dev)
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("WebSocket upgrade error: %v", err)
	}
	defer ws.Close()
	Clients[ws] = true

	for {
		// รออ่าน แต่เราไม่ต้องรับจาก client ตอนนี้
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(Clients, ws)
			break
		}
	}
}

func HandleMessages() {
    for {
        msg := <-Broadcast
        for client := range Clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Printf("WebSocket write error: %v", err)
                client.Close()
                delete(Clients, client)
            }
        }
    }
}
