package models

type Player struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
	RoomID   *int64 `json:"roomId"`
	GameID   *int64 `json:"gameId"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Draws    int    `json:"draws"`
}
