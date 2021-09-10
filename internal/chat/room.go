// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"encoding/json"

	"seabattle/internal/domain"
	"seabattle/internal/game"
)

// Room maintains the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	shot chan *domain.Shoot

	Player1 *Client
	Player2 *Client

	game *game.Service

	chatID string
}

func NewRoom(game *game.Service, chatID string) *Room {
	return &Room{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		shot:       make(chan *domain.Shoot),
		game:       game,
		chatID:     chatID,
	}
}

func (h *Room) Run() {
	for {
		select {
		case client := <-h.register:
			if h.Player1 == nil {
				h.Player1 = client
			} else {
				h.Player2 = client
			}

			h.game.NewPlayer(h.chatID, client.userID)
			gameForUser := h.game.GetGameForUser(h.chatID, client.userID)
			bytes, _ := json.Marshal(gameForUser)
			client.send <- bytes
		case client := <-h.unregister:
			if h.Player1 != nil && h.Player1.userID == client.userID {
				close(h.Player1.send)
				h.Player1 = nil
			}

			if h.Player2 != nil && h.Player2.userID == client.userID {
				close(h.Player2.send)
				h.Player2 = nil
			}
		case shot := <-h.shot:
			if ok := h.game.Shot(h.chatID, shot); ok {
				if h.Player1 != nil {
					gameForUser := h.game.GetGameForUser(h.chatID, h.Player1.userID)
					bytes, _ := json.Marshal(gameForUser)
					h.Player1.send <- bytes
				}

				if h.Player2 != nil {
					gameForUser := h.game.GetGameForUser(h.chatID, h.Player2.userID)
					bytes, _ := json.Marshal(gameForUser)
					h.Player2.send <- bytes
				}
			}
		}
	}
}
