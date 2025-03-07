package tender

import (
	"TenderServiceApi/internal/model"
)

type row struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ServiceType string `db:"service_type"`
	Status      string `db:"status"`
	Version     int    `db:"version"`
	Responsible string `db:"responsible"`
}

func (r *row) toModel() model.Tender {
	return model.Tender{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		ServiceType: r.ServiceType,
		Status:      r.Status,
		Version:     r.Version,
		Responsible: r.Responsible,
	}
}

func toRow(m model.Tender) row {
	return row{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		ServiceType: m.ServiceType,
		Status:      m.Status,
		Version:     m.Version,
		Responsible: m.Responsible,
	}
}
