package model

import (
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

func (t *Tender) FillDefault(tender Tender) {
	if t.Id == "" {
		t.Id = tender.Id
	}

	if t.Name == "" {
		t.Name = tender.Name
	}

	if t.Description == "" {
		t.Description = tender.Description
	}

	if t.ServiceType == "" {
		t.ServiceType = tender.ServiceType
	}

	if t.Status == "" {
		t.Status = tender.Status
	}

	if t.Version == 0 {
		t.Version = tender.Version
	}

	if t.Responsible == "" {
		t.Responsible = tender.Responsible
	}
}

func (a *Tender) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
