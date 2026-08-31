package websocket

import (
"sync"
)

type Hub struct {
    clients    map[int64]*Client 
    broadcast  chan *Message     
    register   chan *Client
    unregister chan *Client
    mu         sync.Mutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[int64]*Client),
        broadcast:  make(chan *Message),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.userID] = client
            h.mu.Unlock()
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.userID]; ok {
                delete(h.clients, client.userID)
                close(client.send)
            }
            h.mu.Unlock()
        case message := <-h.broadcast:

            h.mu.Lock()
            if client, ok := h.clients[message.ReceiverID]; ok {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client.userID)
                }
            }
            h.mu.Unlock()
        }
    }
}