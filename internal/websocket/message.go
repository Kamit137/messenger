package websocket

type Message struct {
    SenderID   int64  `json:"sender_id"`
    ReceiverID int64  `json:"receiver_id"`
    Content    string `json:"content"`
    Type       string `json:"type"` // например, "text", "image"
}