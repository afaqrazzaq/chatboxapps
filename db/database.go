package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ChatRoom struct {
	gorm.Model
	Name  string `json:"name"`                                    //name of users
	Users []User `json:"users" gorm:"many2many:user_chat_rooms;"` // list of users
}

type User struct {
	gorm.Model              // many to many relationship with user and chatroom, chat room can have many users and a user can belong to many chat rooms.
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	ChatRooms    []ChatRoom `json:"chat_rooms" gorm:"many2many:user_chat_rooms;"` //many-to-many relationship between User and ChatRoom
}

type Message struct {
	Messages   []Message `json:"messages"` //Message struct has access to all the fields of the User struct.
	gorm.Model           //
	Text       string    `json:"text"`
	Sender     User      `json:"sender"`
	SenderID   uint      `json:"-"` //foreign id  from user and chatroom
	Room       ChatRoom
	RoomID     uint `json:"-"` //foreign id from user and chatroom
}

func InitDB() (*gorm.DB, error) { // initalizing the connect of mysql database.
	db, err := gorm.Open("mysql", "user:password@tcp(host:port)/database_name?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&ChatRoom{}, &User{}, &Message{})
	return db, nil
}
