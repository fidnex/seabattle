package domain

type State string

const (
	Win   = State("win")
	Loose = State("loose")
	You   = State("you")
	Enemy = State("enemy")
)

const (
	Nothing = 0
	Missed  = 1
	Hit     = 2
	Ship    = 3
)

type Game struct {
	Winner     string
	UserIDTurn string
	ChatID     string
	Player1    Player
	Player2    Player
}

type UserGame struct {
	State State       `json:"state"`
	You   [10][10]int `json:"you"`
	Enemy [10][10]int `json:"enemy"`
}

type Player struct {
	UserID string
	Map    [10][10]int
}

type Shoot struct {
	UserID string `json:"-"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
