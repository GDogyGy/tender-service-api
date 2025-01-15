package tender

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"TenderServiceApi/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (t *Repository) FetchList(ctx context.Context, serviceType string) ([]model.Tender, error) {
	const op = "repository.tender.FetchList"
	var tenders []model.Tender
	var rows *sqlx.Rows
	var err error

	if len(serviceType) > 0 {
		arg := pq.Array(strings.Split(serviceType, ","))
		q := fmt.Sprintf(`SELECT %s FROM tender WHERE tender.service_type = ANY ($1)`, strings.Join(column, ","))
		rows, err = t.db.QueryxContext(ctx, q, arg)
	}

	if len(serviceType) == 0 {
		q := fmt.Sprintf(`SELECT %s FROM tender`, strings.Join(column, ","))
		rows, err = t.db.QueryxContext(ctx, q)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return tenders, model.NotFound
	}

	if err != nil {
		return tenders, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		tender, err := t.tenderFromRows(rows)
		if err != nil {
			return tenders, fmt.Errorf("%s:%w", op, err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (t *Repository) FetchListByUser(ctx context.Context, username string) ([]model.Tender, error) {
	const op = "repository.tender.FetchListByUser"
	var tenders []model.Tender
	var rows *sqlx.Rows
	var err error

	q := fmt.Sprintf(`SELECT %s FROM tender left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1`, strings.Join(column, ","))
	rows, err = t.db.QueryxContext(ctx, q, username)

	if errors.Is(err, sql.ErrNoRows) {
		return tenders, model.NotFound
	}

	if err != nil {
		return tenders, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		tender, err := t.tenderFromRows(rows)
		if err != nil {
			return tenders, fmt.Errorf("%s:%w", op, err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (t *Repository) CheckResponsibleTender(ctx context.Context, username string, tenderId string) (bool, error) {
	const op = "repository.tender.CheckResponsibleTender"
	tender := t.db.QueryRowxContext(ctx, `SELECT COUNT(*) FROM tender left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1 AND tender.id = $2`, username, tenderId)

	err := tender.Err()

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	c := 0

	err = tender.Scan(&c)

	if errors.Is(err, sql.ErrNoRows) {
		return false, model.NotFound
	}

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if c > 0 {
		return true, nil
	}
	return false, model.NotFindResponsibleTender
}

func (t *Repository) FetchTenderById(ctx context.Context, tenderId string) (model.Tender, error) {
	const op = "repository.tender.FetchTenderById"

	q := fmt.Sprintf(`SELECT %s FROM tender WHERE tender.id = $1`, strings.Join(column, ","))
	tender := t.db.QueryRowxContext(ctx, q, tenderId)

	err := tender.Err()

	if errors.Is(err, sql.ErrNoRows) {
		return model.Tender{}, model.NotFound
	}
	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	var te row
	err = tender.StructScan(&te)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Tender{}, model.NotFound
	}

	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	return te.toModel(), nil
}

func (t *Repository) CreateTender(ctx context.Context, saveModel model.Tender) (model.Tender, error) {
	const op = "repository.tender.CreateTender"
	var r row

	r = toRow(saveModel)

	q := "INSERT INTO tender (name, description, service_type, status, responsible) VALUES($1,$2,$3,$4,$5) RETURNING id"

	result := t.db.QueryRowxContext(ctx, q, r.Name, r.Description, r.ServiceType, r.Status, r.Responsible)

	err := result.Err()

	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	var id string

	err = result.Scan(&id)
	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	tender, err := t.FetchTenderById(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Tender{}, model.NotFound
	}

	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	return tender, nil
}

var column = []string{"tender.id", "tender.name", "tender.description", "tender.service_type", "tender.status", "tender.responsible"}

func (t *Repository) tenderFromRows(rows *sqlx.Rows) (model.Tender, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
