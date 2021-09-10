package game

import (
	"seabattle/internal/domain"
	"seabattle/internal/repo"
)

func New(repo *repo.Repo) *Service {
	return &Service{repo: repo}
}

type Service struct {
	repo *repo.Repo
}

func (s *Service) GetGameForUser(chatID, userID string) *domain.UserGame {
	internalGame, _ := s.repo.GetGame(chatID)

	if internalGame == nil {
		internalGame = NewGame(chatID, userID)
		_ = s.repo.SaveGame(internalGame)
	}

	return &domain.UserGame{
		State: s.getGameStateForUser(internalGame, userID),
		You:   s.getOurMap(s.findOurMap(internalGame, userID)),
		Enemy: s.getEnemyMap(s.findEnemyMap(internalGame, userID)),
	}
}

func (s *Service) Shot(chatID string, shoot *domain.Shoot) {
	internalGame, _ := s.repo.GetGame(chatID)
	field := s.findEnemyMap(internalGame, shoot.UserID)

	enemyID := s.getEnemyID(internalGame, shoot.UserID)

	switch field[shoot.X][shoot.Y] {
	case domain.Nothing:
		field[shoot.X][shoot.Y] = domain.Missed
		internalGame.UserIDTurn = enemyID
	case domain.Ship:
		field[shoot.X][shoot.Y] = domain.Hit
	}

}

func (s *Service) findOurMap(game *domain.Game, userID string) [10][10]int {
	if game.Player1.UserID == userID {
		return game.Player1.Map
	}

	return game.Player2.Map
}

func (s *Service) findEnemyMap(game *domain.Game, userID string) [10][10]int {
	if game.Player1.UserID != userID {
		return game.Player2.Map
	}

	return game.Player1.Map
}

func (s *Service) setUserField(game *domain.Game, userID string, field [10][10]int) {
	if game.Player1.UserID == userID {
		game.Player1.Map = field
	}

	game.Player2.Map = field
}

func (s *Service) getGameStateForUser(game *domain.Game, userID string) domain.State {
	switch true {
	case game.Winner != "" && game.Winner == userID:
		return domain.Win
	case game.Winner != "" && game.Winner != userID:
		return domain.Loose
	case game.UserIDTurn == userID:
		return domain.You
	default:
		return domain.Enemy
	}
}

func (s *Service) getOurMap(field [10][10]int) [10][10]int {
	return field
}

func (s *Service) getEnemyMap(field [10][10]int) [10][10]int {
	// меняем кораблики на нолики
	for i := range field {
		for j := range field[i] {
			if field[i][j] == domain.Ship {
				field[i][j] = domain.Nothing
			}
		}
	}

	return field
}

func (s *Service) getEnemyID(game *domain.Game, userID string) string {
	if game.Player1.UserID == userID {
		return game.Player2.UserID
	}

	return game.Player1.UserID
}
