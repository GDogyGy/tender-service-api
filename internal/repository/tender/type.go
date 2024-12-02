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
