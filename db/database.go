package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ChatRoom struct {
	gorm.Model
	Name  string `json:"name"`
	Users []User `json:"users" gorm:"many2many:user_chat_rooms;"`
}

type User struct {
	gorm.Model
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	ChatRooms    []ChatRoom `json:"chat_rooms" gorm:"many2many:user_chat_rooms;"`
	Messages     []Message  `json:"messages"`
}

type Message struct {
	gorm.Model
	Text     string `json:"text"`
	Sender   User   `json:"sender"`
	SenderID uint   `json:"-"`
	Room     ChatRoom
	RoomID   uint `json:"-"`
}

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "user:password@tcp(host:port)/database_name?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&ChatRoom{}, &User{}, &Message{})
	return db, nil
}
