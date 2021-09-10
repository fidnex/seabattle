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
	Winner     string
	UserIDTurn string
	ChatID     string
	Player1    Player
	Player2    Player
}

type UserGame struct {
	State State                   `json:"state"`
	You   [AreaSize][AreaSize]int `json:"you"`
	Enemy [AreaSize][AreaSize]int `json:"enemy"`
}

type Player struct {
	UserID string
	Map    [AreaSize][AreaSize]int
}

type Shoot struct {
	UserID string `json:"-"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
