package organization

import (
	"TenderServiceApi/internal/model"
)

type row struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Type        string `db:"type"`
}

func (r *row) toModel() model.Organization {
	return model.Organization{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Type:        r.Type,
	}
}

func toRow(m model.Organization) row {
	return row{
		Name:        m.Name,
		Description: m.Description,
		Type:        m.Type,
	}
}
