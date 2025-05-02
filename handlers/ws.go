package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"realTimeService/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg models.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		log.Printf("Got message: %+v", msg)
		
		response := map[string]string{
			"status": "ok",
			"echo":   msg.Content,
		}
		resJSON, _ := json.Marshal(response)
		err = conn.WriteMessage(websocket.TextMessage, resJSON)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
