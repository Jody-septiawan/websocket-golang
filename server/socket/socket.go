package socket

import (
	"dumbmerch/models"
	"fmt"
	"log"
	"net/http"

	database "dumbmerch/pkg/mysql"

	"github.com/gorilla/websocket"
)

type PayloadEvent struct {
	Event string	`json:"event"`
}

type PayloadContact struct {
	Event string		`json:"event"`
	Data []models.User	`json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return  true
   },
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	reader(ws)
}

// var db *gorm.DB
func reader(conn *websocket.Conn) {
    for {
		var payloadEvent PayloadEvent
		var payloadContact PayloadContact
		err := conn.ReadJSON(&payloadEvent)
        if err != nil {
            log.Println(err)
            return
        }
		
		if payloadEvent.Event == "customer contacts" {
			var users []models.User
			var chats []models.Chat

			// var err error
			// err := database.DB.Preload("Profile").Preload("Sender").Preload("Recipient").Find(&users, "status = ?", "customer").Error
			err := database.DB.Preload("Profile").Preload("Sender").Find(&users, "status = ?", "customer").Error
			if err != nil {
				log.Println(err)
				return
			}

			err = database.DB.Preload("Sender").Preload("Recipient").Find(&chats).Error
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Println(chats)

			payloadContact.Event = payloadEvent.Event
			payloadContact.Data = users

			conn.WriteJSON(payloadContact)
		} else if payloadEvent.Event == "get admin contact" {
			fmt.Println("========== get admin contact ==========")
		}

        // if err := conn.WriteMessage(messageType, p); err != nil {
        //     log.Println(err)
        //     return
        // }
    }
}