package models

type ChatMessage struct {
	ID        int64  `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
}
