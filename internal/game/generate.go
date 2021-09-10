package game

import "seabattle/internal/domain"

func NewGame(chatID string, userID string) *domain.Game {
	return &domain.Game{
		ChatID: chatID,
		Player1: domain.Player{
			UserID: userID,
			Map:    generateMap(),
		},
		Player2: domain.Player{
			UserID: "", // черт кто знает с кем играем
			Map:    generateMap(),
		},
	}
}

func generateMap() [10][10]int {
	return [10][10]int{
		{0, 3, 0, 3, 0, 3, 3, 0, 0, 0},
		{0, 3, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 3, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 3, 0, 0, 0, 3, 3, 3, 0, 0},
		{0, 0, 0, 3, 0, 0, 0, 0, 0, 0},
		{3, 0, 0, 0, 0, 0, 0, 0, 3, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 3, 0},
		{0, 0, 0, 3, 3, 3, 0, 0, 0, 0},
		{0, 3, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
	}
}
