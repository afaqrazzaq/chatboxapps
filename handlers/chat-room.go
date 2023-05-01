    // Routes for message handlers
    r.HandleFunc("/chatrooms/{id}/messages", handlers.SendMessage).Methods("POST")
    r.HandleFunc("/chatrooms/{id}/messages", handlers.GetMessages).Methods("GET")

package handlers

import (  // provides support for WebSocket connections
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{    //used to upgrade a standard HTTP connection to a WebSocket connection.
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// Listen for messages from the client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// Print the received message to the console
		fmt.Printf("Received message: %s\n", p)

		// Write the message back to the client
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println("Error writing message:", err)
			return
		}
	}
}
// CreateChatRoom handles POST requests to create a new chat room
func CreateChatRoom(w http.ResponseWriter, r *http.Request) {
    // Parse request body to get chat room name
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    
    // Create chat room in database
    chatRoomName := string(body)
    chatRoom, err := models.CreateChatRoom(chatRoomName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Send response with chat room information
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(chatRoom)
}

// JoinChatRoom handles POST requests to join an existing chat room
func JoinChatRoom(w http.ResponseWriter, r *http.Request) {
    // Parse request body to get user and chat room IDs
    var reqBody struct {
        UserID    int64 `json:"userId"`
        ChatRoomID int64 `json:"chatRoomId"`
    }
    err := json.NewDecoder(r.Body).Decode(&reqBody)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Join user to chat room in database
    err = models.JoinChatRoom(reqBody.UserID, reqBody.ChatRoomID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Send success response
    w.WriteHeader(http.StatusOK)
}

// LeaveChatRoom handles POST requests to leave an existing chat room
func LeaveChatRoom(w http.ResponseWriter, r *http.Request) {
    // Parse request body to get user and chat room IDs
    var reqBody struct {
        UserID    int64 `json:"userId"`
        ChatRoomID int64 `json:"chatRoomId"`
    }
    err := json.NewDecoder(r.Body).Decode(&reqBody)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Leave user from chat room in database
    err = models.LeaveChatRoom(reqBody.UserID, reqBody.ChatRoomID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Send success response
    w.WriteHeader(http.StatusOK)
}
// UserConnected handles a new user connection to the chat room.
func (cr *ChatRoom) UserConnected(conn *websocket.Conn) {
    userID, err := cr.authenticateUser(conn)
    if err != nil {
        log.Printf("Failed to authenticate user: %v", err)
        return
    }

    cr.users[userID] = conn

    // Send a welcome message to the new user.
    welcomeMsg := fmt.Sprintf("Welcome, %v!", userID)
    err = conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))
    if err != nil {
        log.Printf("Failed to send welcome message: %v", err)
    }

    // Broadcast user connection to all connected users.
    cr.broadcastUserConnected(userID)
}

// UserDisconnected handles a user disconnection from the chat room.
func (cr *ChatRoom) UserDisconnected(userID string) {
    delete(cr.users, userID)

    // Broadcast user disconnection to all connected users.
    cr.broadcastUserDisconnected(userID)
}
