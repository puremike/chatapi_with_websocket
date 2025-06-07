package ws

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomID]; ok {
				if _, ok := h.Rooms[client.RoomID].Clients[client.ID]; !ok {
					h.Rooms[client.RoomID].Clients[client.ID] = client
				}
			}
		case unRegClient := <-h.Unregister:
			if _, ok := h.Rooms[unRegClient.RoomID]; ok {
				if _, ok := h.Rooms[unRegClient.RoomID].Clients[unRegClient.ID]; ok {

					if len(h.Rooms[unRegClient.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  unRegClient.Username + " has left the room.",
							RoomID:   unRegClient.RoomID,
							Username: unRegClient.Username,
						}
					}
					delete(h.Rooms[unRegClient.RoomID].Clients, unRegClient.ID)
					close(unRegClient.Message)
				}
			}
		case msg := <-h.Broadcast:
			if _, ok := h.Rooms[msg.RoomID]; ok {
				for _, client := range h.Rooms[msg.RoomID].Clients {
					client.Message <- msg
				}
			}
		}
	}
}
