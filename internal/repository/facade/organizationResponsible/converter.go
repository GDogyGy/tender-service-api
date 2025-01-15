package organizationResponsible

import "TenderServiceApi/internal/model"

type row struct {
	Id             string `db:"id"`
	OrganizationId string `db:"organization_id"`
	UserId         string `db:"user_id"`
}

func (r *row) toModel() model.OrganizationResponsible {
	return model.OrganizationResponsible{
		Id:             r.Id,
		OrganizationId: r.OrganizationId,
		UserId:         r.UserId,
	}
}
