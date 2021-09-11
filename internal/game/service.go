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

func (s *Service) GetGame(chatID, userID string) *domain.Game {
	internalGame, _ := s.repo.GetGame(chatID)

	if internalGame == nil {
		internalGame = NewGame(chatID, userID)
		_ = s.repo.SaveGame(internalGame)
	}

	return internalGame
}

func (s *Service) GetGameForUser(chatID, userID string) *domain.UserGame {
	internalGame := s.GetGame(chatID, userID)

	return &domain.UserGame{
		State: s.getGameStateForUser(internalGame, userID),
		You:   s.getOurMap(s.findOurMap(internalGame, userID)),
		Enemy: s.getEnemyMap(s.findEnemyMap(internalGame, userID)),
	}
}

func (s *Service) Shot(chatID string, shoot *domain.Shoot) bool {
	internalGame, _ := s.repo.GetGame(chatID)
	if internalGame.UserIDTurn != shoot.UserID {
		return false
	}

	field := s.findEnemyMap(internalGame, shoot.UserID)

	enemyID := s.getEnemyID(internalGame, shoot.UserID)

	switch field[shoot.X][shoot.Y] {
	case domain.Nothing:
		field[shoot.X][shoot.Y] = domain.Missed
		internalGame.UserIDTurn = enemyID
	case domain.Ship:
		field[shoot.X][shoot.Y] = domain.Hit
		if ok := s.shipIsDrowned(field, shoot); ok {
			field = s.drownShip(field, shoot)
		}
	default:
		return false
	}

	s.setUserField(internalGame, enemyID, field)

	if win := s.isUserWin(field); win {
		internalGame.Winner = shoot.UserID
	}

	_ = s.repo.SaveGame(internalGame)

	return true
}

func (s *Service) NewPlayer(chatID, userID string) {
	game := s.GetGame(chatID, userID)

	if game.Player1.UserID == userID {
		return
	}

	if game.Player2.UserID == userID {
		return
	}

	if game.Player1.UserID == "" {
		game.Player1.UserID = userID
	}

	if game.Player2.UserID == "" {
		game.Player2.UserID = userID
	}

	if game.UserIDTurn == "" {
		game.UserIDTurn = userID
	}

	_ = s.repo.SaveGame(game)
}

func (s *Service) findOurMap(game *domain.Game, userID string) [10][10]int {
	if game.Player1.UserID == userID {
		return game.Player1.Map
	}

	return game.Player2.Map
}

func (s *Service) findEnemyMap(game *domain.Game, userID string) [10][10]int {
	if game.Player1.UserID == userID {
		return game.Player2.Map
	}

	return game.Player1.Map
}

func (s *Service) setUserField(game *domain.Game, userID string, field [10][10]int) {
	if game.Player1.UserID == userID {
		game.Player1.Map = field
		return
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

func (s *Service) shipIsDrowned(field [10][10]int, shoot *domain.Shoot) bool {
	var i = 0
	iterateShip := true
	for iterateShip {
		i++
		iterateShip = false
		if shoot.X-i >= 0 {
			val := field[shoot.X-i][shoot.Y]

			if val == domain.Ship {
				return false
			}

			if val == domain.Hit {
				iterateShip = true
			}
		}
	}

	i = 0
	iterateShip = true
	for iterateShip {
		i++
		iterateShip = false
		if shoot.X+i <= 9 {
			val := field[shoot.X+i][shoot.Y]

			if val == domain.Ship {
				return false
			}

			if val == domain.Hit {
				iterateShip = true
			}
		}
	}

	i = 0
	iterateShip = true
	for iterateShip {
		i++
		iterateShip = false
		if shoot.Y-i >= 0 {
			val := field[shoot.X][shoot.Y-i]

			if val == domain.Ship {
				return false
			}

			if val == domain.Hit {
				iterateShip = true
			}
		}
	}

	i = 0
	iterateShip = true
	for iterateShip {
		i++
		iterateShip = false
		if shoot.Y+i <= 9 {
			val := field[shoot.X][shoot.Y+i]
			if val == domain.Ship {
				return false
			}

			if val == domain.Hit {
				iterateShip = true
			}
		}
	}

	return true
}

func (s *Service) drownShip(field [10][10]int, shoot *domain.Shoot) [10][10]int {
	isHorizontal := func() bool {
		val := shoot.Y - 1

		if val >= 0 && field[shoot.X][val] == domain.Hit {
			return true
		}

		val = shoot.Y + 1
		if val <= 9 && field[shoot.X][val] == domain.Hit {
			return true
		}
		return false
	}()

	isVertical := func() bool {
		val := shoot.X - 1

		if val >= 0 && field[val][shoot.Y] == domain.Hit {
			return true
		}

		val = shoot.X + 1
		if val <= 9 && field[val][shoot.Y] == domain.Hit {
			return true
		}

		return false
	}()

	fillCell := func(x, y int) {
		if x > 9 || x < 0 {
			return
		}

		if y > 9 || y < 0 {
			return
		}

		if field[x][y] == domain.Nothing {
			field[x][y] = domain.Missed
		}
	}

	outline := func(x, y int) {
		fillCell(x+1, y-1)
		fillCell(x+1, y)
		fillCell(x+1, y+1)

		fillCell(x, y-1)
		fillCell(x, y+1)

		fillCell(x-1, y-1)
		fillCell(x-1, y)
		fillCell(x-1, y+1)
	}

	outline(shoot.X, shoot.Y)

	// вертикальный кораблик
	if isVertical {
		i := 0
		for {
			i++
			val := shoot.X + i
			if val <= 9 {
				if field[val][shoot.Y] == domain.Hit {
					outline(val, shoot.Y)
				} else {
					break
				}
			} else {
				break
			}
		}

		i = 0
		for {
			i++
			val := shoot.X - i
			if val >= 0 {
				if field[val][shoot.Y] == domain.Hit {
					outline(val, shoot.Y)
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	// горизонтальный кораблик
	if isHorizontal {
		i := 0
		for {
			i++
			val := shoot.Y + i
			if val <= 9 {
				if field[shoot.X][val] == domain.Hit {
					outline(shoot.X, val)
				} else {
					break
				}
			} else {
				break
			}
		}

		i = 0
		for {
			i++
			val := shoot.Y - i
			if val >= 0 {
				if field[shoot.X][val] == domain.Hit {
					outline(shoot.X, val)
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	return field
}

func (s *Service) isUserWin(field [10][10]int) bool {
	for i := range field {
		for j := range field[i] {
			if field[i][j] == domain.Ship {
				return false
			}
		}
	}

	return true
}
