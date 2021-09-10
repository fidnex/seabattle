package game

import (
	"crypto/rand"
	"math"
	"math/big"

	"seabattle/internal/domain"
)

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

// generateMap генерация карты боя
func generateMap() [domain.AreaSize][domain.AreaSize]int {
	area := [domain.AreaSize][domain.AreaSize]int{}

	// 1 корабль - ряд из 4 клеток («четырёхпалубный»; линкор)
	generateShip(4, &area)

	// 2 корабля - ряд из 3 клеток («трёхпалубные»; крейсера)
	for i := 0; i < 2; i++ {
		generateShip(3, &area)
	}

	// 3 корабля - ряд из 2 клеток («двухпалубные»; эсминцы)
	for i := 0; i < 3; i++ {
		generateShip(2, &area)
	}

	// 4 корабля - 1 клетка («однопалубные»; торпедные катера)
	for i := 0; i < 4; i++ {
		generateShip(1, &area)
	}

	return area
}

// generateShip генерируем расположение корабля заданной длины.
func generateShip(shipSize int, area *[domain.AreaSize][domain.AreaSize]int) {
	var (
		x, y           int  // Начальные координаты корабля
		isHorizontally bool // Горизонтальное расположение корабля
	)

	// Рандомно определяем ориентацию корабля (по горизонтали или по вертикали)
	rnd, _ := rand.Int(rand.Reader, big.NewInt(2))
	if (rnd.Int64() % 2) == 0 {
		isHorizontally = true
	}

	// Тупо псевдослучано накидываем корабли на карту :)
	for {
		for {
			rnd, _ := rand.Int(rand.Reader, big.NewInt(domain.AreaSize))
			x = int(rnd.Int64()) % domain.AreaSize

			if !(!isHorizontally && x > domain.AreaSize-shipSize) {
				break
			}
		}

		for {
			rnd, _ := rand.Int(rand.Reader, big.NewInt(domain.AreaSize))
			y = int(rnd.Int64()) % domain.AreaSize

			if !(isHorizontally && y > domain.AreaSize-shipSize) {
				break
			}
		}

		if validateShip(shipSize, isHorizontally, x, y, area) {
			break
		}
	}

	// Отмечаем корабль на карте
	if isHorizontally {
		for j := y; j < y+shipSize; j++ {
			area[x][j] = domain.Ship
		}
	} else {
		for i := x; i < x+shipSize; i++ {
			area[i][y] = domain.Ship
		}
	}
}

// validateShip проверяем корабль на карте на предмет адекватности.
func validateShip(shipSize int, isHorizontally bool, x, y int, area *[domain.AreaSize][domain.AreaSize]int) bool {
	if isHorizontally {
		for i := int(math.Max(0, float64(x-1))); i <= int(math.Min(domain.AreaSize-1, float64(x+1))); i++ {
			for j := int(math.Max(0, float64(y-1))); j <= int(math.Min(domain.AreaSize-1, float64(y+shipSize))); j++ {
				if area[i][j] == domain.Ship {
					return false
				}
			}
		}
	} else {
		for i := int(math.Max(0, float64(x-1))); i <= int(math.Min(domain.AreaSize-1, float64(x+shipSize))); i++ {
			for j := int(math.Max(0, float64(y-1))); j <= int(math.Min(domain.AreaSize-1, float64(y+1))); j++ {
				if area[i][j] == domain.Ship {
					return false
				}
			}
		}
	}

	return true
}
