package models

import "time"

type User struct {
	UUID         string    `json:"uuid"`
	PhoneNumber  string    `json:"phoneNumber"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	RegisteredAt time.Time `json:"registeredAt"`
	IsAdmin      bool      `json:"isAdmin"`
}
