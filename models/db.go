package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return &DB{db}, nil
}

func (db *DB) InsertChatMessage(chatMessage *ChatMessage) (int64, error) {
	query := `INSERT INTO chat_messages (sender, recipient, message, created_at) 
                VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err := db.QueryRow(query, chatMessage.Sender, chatMessage.Recipient, chatMessage.Message, chatMessage.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) GetChatMessages(sender, recipient string) ([]*ChatMessage, error) {
	query := `SELECT * FROM chat_messages WHERE (sender=$1 AND recipient=$2) OR (sender=$2 AND recipient=$1) ORDER BY created_at ASC`

	rows, err := db.Query(query, sender, recipient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatMessages := []*ChatMessage{}
	for rows.Next() {
		chatMessage := &ChatMessage{}
		err := rows.Scan(&chatMessage.ID, &chatMessage.Sender, &chatMessage.Recipient, &chatMessage.Message, &chatMessage.CreatedAt)
		if err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatMessages, nil
}
