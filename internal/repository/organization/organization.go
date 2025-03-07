package organization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"TenderServiceApi/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (o *Repository) FetchById(ctx context.Context, id string) (model.Organization, error) {
	const op = "repository.organization.FetchById"
	var organization model.Organization
	if len(id) == 0 {
		return organization, fmt.Errorf("%s:%s", op, "Parametr id is empty")
	}

	r := o.db.QueryRowx(fmt.Sprintf(`SELECT %s FROM organization WHERE id = $1`, strings.Join(column, ", ")), id)

	organizationRow, err := o.organizationFromRow(r)
	if errors.Is(err, sql.ErrNoRows) {
		return organization, fmt.Errorf("%s:%w", op, err)
	}
	if err != nil {
		return organization, err
	}

	return organizationRow, nil
}

func (o *Repository) CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error) {
	const op = "repository.organization.CheckResponsible"
	var result model.OrganizationResponsible
	q := "SELECT organization_responsible.id, organization_responsible.organization_id, organization_responsible.user_id FROM organization_responsible left join employee on employee.id = organization_responsible.user_id left join organization on organization.id = organization_responsible.organization_id WHERE employee.username = $1 AND organization.id = $2"

	row := o.db.QueryRowxContext(ctx, q, username, organizationId)
	err := row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return result, model.NotFound
	}
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	result, err = o.organizationResponsibleFromRow(row)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (o *Repository) FetchRelationsById(ctx context.Context, id string) (model.OrganizationResponsible, error) {
	const op = "repository.organization.FetchRelationsById"
	var result model.OrganizationResponsible

	q := "SELECT organization_responsible.id, organization_responsible.organization_id, organization_responsible.user_id FROM organization_responsible WHERE organization_responsible.id = $1"

	row := o.db.QueryRowxContext(ctx, q, id)
	err := row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return result, model.NotFound
	}
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	result, err = o.organizationResponsibleFromRow(row)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

var column = []string{"organization.id", "organization.name", "organization.description", "organization.type"}

func (e *Repository) organizationFromRow(ro *sqlx.Row) (model.Organization, error) {
	var r organizationRow
	err := ro.StructScan(&r)
	m := r.toModel()
	return m, err
}

func (e *Repository) organizationResponsibleFromRow(ro *sqlx.Row) (model.OrganizationResponsible, error) {
	var r organizationResponsibleRow
	err := ro.StructScan(&r)
	m := r.toModel()
	return m, err
}
