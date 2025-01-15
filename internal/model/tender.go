package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Tender struct {
	Id          string
	Name        string
	Description string
	ServiceType string
	Status      string
	Version     int
	Responsible string
}

func (a Tender) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Tender) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
