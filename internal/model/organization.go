package model

import "time"

type Organization struct {
	Id          string
	Name        string
	Description string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
