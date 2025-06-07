package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		msg, ok := <-c.Message
		if !ok {
			return
		}
		if err := c.Conn.WriteJSON(msg); err != nil {
			return
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		hub.Broadcast <- msg
	}
}
