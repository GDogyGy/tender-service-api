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

// TODO: Обсудить с димой замену рефлекту Возможно нужно проще if (id == "") d.Id = a.Id
func (a *Tender) FillDefault(defaults *Tender) {
	type field struct {
		src  interface{}
		dest interface{}
	}

	fields := []field{
		{&a.Id, &defaults.Id},
		{&a.Name, &defaults.Name},
		{&a.Description, &defaults.Description},
		{&a.ServiceType, &defaults.ServiceType},
		{&a.Status, &defaults.Status},
		{&a.Version, &defaults.Version},
		{&a.Responsible, &defaults.Responsible},
	}

	for _, f := range fields {
		switch dest := f.dest.(type) {
		case *string:
			if *dest == "" {
				*dest = *f.src.(*string)
			}
		case *int:
			if *dest == 0 {
				*dest = *f.src.(*int)
			}
		}
	}
}

func (a *Tender) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
