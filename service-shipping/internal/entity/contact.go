package entity

type Contact struct {
	Id     int16  `json:"id"`
	UserId string `json:"user_id"`
	ContactInfo
	CreateDelete
}
