package models

type RoomStatus int

const (
	ROOM_OPEN   RoomStatus = 1
	ROOM_FULL   RoomStatus = 2
	ROOM_CLOSED RoomStatus = 3
)

type RoomParticipant struct {
	ID       *int64 `json:"id"`
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
	GameID      *int64          `json:"gameId"`
	PrevGameID  *int64          `json:"prevGameId"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      RoomStatus      `json:"status"`
}
