package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *WSHandler) CreateRoom(c *gin.Context) {
	var payload CreateRoomRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[payload.ID] = &Room{
		ID:      payload.ID,
		Name:    payload.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusCreated, gin.H{"message": "room created", "room_id": payload.ID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity, adjust as needed
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *WSHandler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}
	defer conn.Close()

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		Username: username,
		RoomID:   roomID,
	}

	m := &Message{
		Content:  username + " has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m
	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *WSHandler) GetRooms(c *gin.Context) {

	rooms := []RoomRes{}

	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *WSHandler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomID := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomID]; !ok {
		clients := make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, client := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
