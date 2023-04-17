package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Define a message struct
type Message struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"room_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// A slice of messages
var messages []Message

// GetMessages returns all messages in the given chat room
func GetMessages(w http.ResponseWriter, r *http.Request) {
	// Get the chat room ID from the request URL
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	// Filter the messages slice by the chat room ID
	filteredMessages := make([]Message, 0)
	for _, m := range messages {
		if m.RoomID == roomID {
			filteredMessages = append(filteredMessages, m)
		}
	}

	// Send the filtered messages as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredMessages)
}

// PostMessage creates a new message in the given chat room
func PostMessage(w http.ResponseWriter, r *http.Request) {
	// Get the chat room ID and message content from the request body
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	message.RoomID = roomID

	// Generate a unique ID for the new message
	lastMessage := messages[len(messages)-1]
	message.ID = lastMessage.ID + 1

	// Add the new message to the messages slice
	messages = append(messages, message)

	log.Printf("Created message %+v\n", message)

	// Send the new message as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}


// handlers/message.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Message struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"room_id"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func SendMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg Message
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO messages (room_id, sender_id, content) VALUES ($1, $2, $3)`
		_, err = db.Exec(query, msg.RoomID, msg.SenderID, msg.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

func GetMessages(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room_id")

		query := `SELECT id, room_id, sender_id, content, timestamp FROM messages WHERE room_id = $1 ORDER BY timestamp`
		rows, err := db.Query(query, roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var messages []Message
		for rows.Next() {
			var msg Message
			err := rows.Scan(&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.Timestamp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			messages = append(messages, msg)
		}
		err = rows.Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}
import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    "time"
    "your-app-name/models"
    "your-app-name/utils"
)
