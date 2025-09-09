package models

type GameStatus int

const (
	On GameStatus = 1
)

type GamePlayer struct {
	ID   int64
	Mark rune
}

type Game struct {
	ID              int64
	Host            GamePlayer
	Guest           GamePlayer
	CurrentPlayerID int
	Board           string
	Status          GameStatus
	WinnerID        int
	LoserID         int
}
