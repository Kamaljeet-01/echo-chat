package user

import "gorm.io/gorm"

type Chatuser struct {
	gorm.Model // adds ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string
	Email      string
}

func (t Chatuser) TableName() string { // This sets the table name to "Chat_user" for the Chatuser struct.
	return "chat_user"
}
