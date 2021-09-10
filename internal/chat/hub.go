package chat

import "seabattle/internal/game"

func NewHub(gameService *game.Service) *Hub {
	return &Hub{rooms: make(map[string]*Room), gameService: gameService}
}

type Hub struct {
	rooms       map[string]*Room
	gameService *game.Service
}

func (h *Hub) GetRoom(chatID string) *Room {
	return h.rooms[chatID]
}

func (h *Hub) CreateNew(chatID string) *Room {
	room := NewRoom(h.gameService, chatID)
	h.rooms[chatID] = room
	return room
}
