package organization

import (
	"TenderServiceApi/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
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
	if len(id) <= 0 {
		return organization, fmt.Errorf("%s:%s", op, "Parametr id is empty")
	}

	r := o.db.QueryRowx(fmt.Sprintf(`SELECT %s FROM organization WHERE id = $1`, strings.Join(column, ", ")), id)

	organizationRow, err := o.fromRow(r)

	if errors.Is(err, sql.ErrNoRows) {
		return organization, fmt.Errorf("%s:%w", op, err)
	}

	if err != nil {
		return organization, err
	}

	return organizationRow, nil
}

var column = []string{"organization.id", "organization.name", "organization.description", "organization.type"}

func (o *Repository) fromRows(rows *sqlx.Rows) (model.Organization, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}

func (e *Repository) fromRow(ro *sqlx.Row) (model.Organization, error) {
	var r row
	err := ro.StructScan(&r)
	m := r.toModel()
	return m, err
}
