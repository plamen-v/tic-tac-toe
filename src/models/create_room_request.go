package models

type CreateRoomRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
