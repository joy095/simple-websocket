package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]struct{})

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin (for dev)
		return true
	},
}

func main() {
	router := gin.Default()

	// Enable CORS for frontend during development
	router.Use(cors.Default())

	router.GET("/ws", serveWs)

	log.Println("Server starting on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

func serveWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	clients[conn] = struct{}{}
	log.Println("Client connected")

	go handleClient(conn)
}

func handleClient(conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
		log.Println("Client disconnected")
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received message: %+v", msg)
		broadcast(msg)
	}
}

func broadcast(msg Message) {
	for conn := range clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Broadcast error: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
