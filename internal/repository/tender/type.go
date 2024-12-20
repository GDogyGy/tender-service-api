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
	Responsible string `db:"responsible"`
}

type createDTO struct {
	Name            string `db:"name"`
	Description     string `db:"description"`
	ServiceType     string `db:"serviceType"`
	Status          string `db:"status"`
	OrganizationId  string `db:"organizationId"`
	CreatorUsername string `db:"creatorUsername"`
}

func (r *row) toModel() model.Tender {
	return model.Tender{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		ServiceType: r.ServiceType,
		Status:      r.Status,
		Responsible: r.Responsible,
	}
}

func toRow(m model.Tender) row {
	return row{
		Name:        m.Name,
		Description: m.Description,
		ServiceType: m.ServiceType,
		Status:      m.Status,
		Responsible: m.Responsible,
	}
}
