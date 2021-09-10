package repo

import (
	"seabattle/internal/domain"
)

type Repo struct {
	games map[string]*domain.Game
}

func New() *Repo {
	return &Repo{games: make(map[string]*domain.Game)}
}

func (r *Repo) GetGame(chatID string) (*domain.Game, error) {
	if game, ok := r.games[chatID]; ok {
		return game, nil
	}

	return nil, nil
}

func (r *Repo) SaveGame(game *domain.Game) error {
	r.games[game.ChatID] = game
	return nil
}
