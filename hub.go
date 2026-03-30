package main

import (
	"context"
	"regexp"

	"nhooyr.io/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	history    [][]byte
}

func NewHub() *Hub {
	x := Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		history:    make([][]byte, 0),
	}
	return &x
}
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	x := Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	return &x
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			for _, m := range h.history {
				client.send <- m
			}
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
		case msg := <-h.broadcast:
			re := regexp.MustCompile(`^.+: .+$`)
			good := re.MatchString(string(msg))
			if good {
				msg = append(msg, byte('\n'))
				h.history = append(h.history, msg)
				for client := range h.clients {
					client.send <- msg
				}
			}
		}

	}
}
func (c *Client) ReadPump(ctx context.Context) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			c.hub.unregister <- c
			return
		}
		c.hub.broadcast <- data
	}
}
func (c *Client) WritePump(ctx context.Context) {
	for {
		select {
		case msg := <-c.send:
			err := c.conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {

			}
		}
	}
}
