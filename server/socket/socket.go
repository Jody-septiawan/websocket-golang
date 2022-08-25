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
	Event 		string	`json:"event"`
	SenderID 	int		`json:"senderId"`
	RecipientID int		`json:"recipientId"`
	Message 	string 	`json:"message"`
}

type PayloadContact struct {
	Event string		`json:"event"`
	Data []models.User	`json:"data"`
}

type PayloadMessage struct {
	Event string		`json:"event"`
	Data []models.Chat	`json:"data"`
}

var clients = make(map[*websocket.Conn]int)

var payloadEvent PayloadEvent
var payloadContact PayloadContact
var payloadMessage PayloadMessage

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

	r.Header.Get("")

	token := r.Header.Get("Authorization")
	
	log.Println("Token : ", token)

	handler(ws)
}

func handler(conn *websocket.Conn) {
	var users []models.User
	var chats []models.Chat
    for {
		err := conn.ReadJSON(&payloadEvent)
        if err != nil {
            log.Println(err)
            return
        }

		fmt.Println(payloadEvent.Event)
		fmt.Println(clients)
		fmt.Println("Client Count : ", len(clients))
		
		if payloadEvent.Event == "customer contacts" {
			saveClient(conn, payloadEvent.SenderID)
			err := database.DB.Preload("Profile").Preload("Chat").Find(&users, "status = ?", "customer").Error
			if err != nil {
				log.Println(err)
				return
			}

			payloadContact.Event = payloadEvent.Event
			payloadContact.Data = users
			
			conn.WriteJSON(payloadContact)
		} else if payloadEvent.Event == "admin contact" {
			saveClient(conn, payloadEvent.SenderID)
			err := database.DB.Preload("Profile").Preload("Chat").Find(&users, "status = ?", "admin").Error
			if err != nil {
				log.Println(err)
				return
			}

			payloadContact.Event = payloadEvent.Event
			payloadContact.Data = users
			
			conn.WriteJSON(payloadContact)
		} else if payloadEvent.Event == "messages" {
			err := database.DB.Debug().Preload("Sender").Preload("Recipient").Find(&chats, "(recipient_id = ? OR recipient_id = ?) AND (sender_id = ? OR sender_id = ?)", payloadEvent.SenderID, payloadEvent.RecipientID, payloadEvent.SenderID, payloadEvent.RecipientID).Error
			if err != nil {
				log.Println(err)
				return
			}
			
			payloadMessage.Event = payloadEvent.Event
			payloadMessage.Data = chats

			conn.WriteJSON(payloadMessage)
		} else if payloadEvent.Event == "send message" {
			var chat models.Chat

			chat.Message = payloadEvent.Message
			chat.SenderID = payloadEvent.SenderID
			chat.RecipientID = payloadEvent.RecipientID

			err := database.DB.Debug().Create(&chat).Error
			if err != nil {
				log.Println(err)
				return
			}
			
			payloadEvent.Event = "new message"
			for client := range clients {
				fmt.Println("client: ", client)
				err :=  client.WriteJSON(payloadEvent)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
    }
}

func saveClient(ws *websocket.Conn, userId int) {
	var already = false
	for _, elm := range clients {
		if elm == userId {
			already = true
		}
	}
	if !already {
		log.Println("Client Connected : ", userId)
		clients[ws] = userId
	}
}