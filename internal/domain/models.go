package domain

type State string

const AreaSize = 10 // Размер матрицы

const (
	Win   = State("win")
	Loose = State("loose")
	You   = State("you")
	Enemy = State("enemy")
)

const (
	Nothing = 0 // непонятно
	Missed  = 1 // был выстрел, промах
	Hit     = 2 // попал в корабль
	Ship    = 3 // корабль (пользователя)
	Drowned = 4 // корабль (врага)
)

type Game struct {
	Winner     string `json:"winner"`
	UserIDTurn string `json:"userIDTurn"`
	ChatID     string `json:"chatID"`
	Player1    Player `json:"player1"`
	Player2    Player `json:"player2"`
}

type UserGame struct {
	State State                   `json:"state"`
	You   [AreaSize][AreaSize]int `json:"you"`
	Enemy [AreaSize][AreaSize]int `json:"enemy"`
}

type Player struct {
	UserID string                  `json:"userID"`
	Map    [AreaSize][AreaSize]int `json:"map"`
}

type Shoot struct {
	UserID string `json:"-"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
