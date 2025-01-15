package model

import "time"

type Employee struct {
	Id        string
	UserName  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
