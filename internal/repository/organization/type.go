package organization

import (
	"TenderServiceApi/internal/model"
)

type organizationRow struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Type        string `db:"type"`
}

type organizationResponsibleRow struct {
	Id             string `db:"id"`
	OrganizationId string `db:"organization_id"`
	UserId         string `db:"user_id"`
}

func (r *organizationRow) toModel() model.Organization {
	return model.Organization{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Type:        r.Type,
	}
}

func toRow(m model.Organization) organizationRow {
	return organizationRow{
		Name:        m.Name,
		Description: m.Description,
		Type:        m.Type,
	}
}

func (o *organizationResponsibleRow) toModel() model.OrganizationResponsible {
	return model.OrganizationResponsible{
		Id:             o.Id,
		OrganizationId: o.OrganizationId,
		UserId:         o.UserId,
	}
}
