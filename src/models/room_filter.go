package models

type RoomFilter struct {
	Host        string     `json:"host"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      RoomStatus `json:"status"`
}
