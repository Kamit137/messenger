package websocket

import (
    "github.com/gorilla/websocket"
    "log"
)

type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    send   chan *Message
    userID int64
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    for {
        var msg Message
        err := c.conn.ReadJSON(&msg)
        if err != nil {
            break
        }
        // Отправляем сообщение в Hub
        msg.SenderID = c.userID
        c.hub.broadcast <- &msg
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()
    for msg := range c.send {
        err := c.conn.WriteJSON(msg)
        if err != nil {
            break
        }
    }
}