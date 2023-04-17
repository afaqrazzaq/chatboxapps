package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
) // Define database connection parameters
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "chatapp"
)

// Establish a connection with the database
func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database")
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Your code here
}
func main() {
	// Initialize the router
	r := mux.NewRouter()

	// Define the chat room routes
	r.HandleFunc("/chat-room/create", handlers.CreateChatRoom).Methods("POST")
	r.HandleFunc("/chat-room/join", handlers.JoinChatRoom).Methods("POST")
	r.HandleFunc("/chat-room/leave", handlers.LeaveChatRoom).Methods("POST")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
