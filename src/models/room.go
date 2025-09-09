package models

type RoomStatus int

const (
	OPEN   RoomStatus = 1
	FULL   RoomStatus = 2
	CLOSED RoomStatus = 3
)

type RoomParticipant struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsReady  bool   `json:"isReady"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Draws    int    `json:"draws"`
}

type Room struct {
	ID          int64           `json:"id"`
	Host        RoomParticipant `json:"host"`
	Guest       RoomParticipant `json:"guest"`
	GameID      int             `json:"gameId"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      RoomStatus      `json:"status"`
}
