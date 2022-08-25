package models

import "time"

type Chat struct {
	ID        	int              	`json:"id" gorm:"primary_key:auto_increment"`
	SenderID  	int                	`json:"sender_id"`
	Sender    	UsersChatResponse	`json:"sender"`
	RecipientID	int                	`json:"recipient_id"`
	Recipient	UsersChatResponse	`json:"recipient"`
	Message		string			   	`json:"message"`
	CreatedAt	time.Time           `json:"-"`
	UpdatedAt	time.Time           `json:"-"`
}

type ChatSenderResponse struct {
	ID        	int              	`json:"id"`
	SenderID  	int                	`json:"-"`
	Message		string			   	`json:"message"`
}

func (ChatSenderResponse) TableName() string {
	return "chats"
}

type ChatRecipientResponse struct {
	ID        	int              	`json:"id"`
	RecipientID  	int             `json:"-"`
	Message		string			   	`json:"message"`
}

func (ChatRecipientResponse) TableName() string {
	return "chats"
}