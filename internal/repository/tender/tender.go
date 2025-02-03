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
		q := fmt.Sprintf(`SELECT %s FROM tender WHERE version = (SELECT MAX(version) FROM tender t2 WHERE tender.id = t2.id and t2.status = 'PUBLISHED' and t2.service_type = ANY ($1))`, strings.Join(column, ","))
		rows, err = t.db.QueryxContext(ctx, q, arg)
	}

	if len(serviceType) == 0 {
		q := fmt.Sprintf("SELECT %s FROM tender WHERE version = (SELECT MAX(version) FROM tender t2 WHERE tender.id = t2.id and t2.status = 'PUBLISHED')", strings.Join(column, ","))
		rows, err = t.db.QueryxContext(ctx, q)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return tenders, model.NotFound
	}

	if err != nil {
		return tenders, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		tender, err := t.fromRows(rows)
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

	q := fmt.Sprintf(`SELECT %s FROM tender left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1 and version = (SELECT MAX(version) FROM tender t2 WHERE tender.id = t2.id)`, strings.Join(column, ","))
	rows, err = t.db.QueryxContext(ctx, q, username)
	if errors.Is(err, sql.ErrNoRows) {
		return tenders, model.NotFound
	}
	if err != nil {
		return tenders, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		tender, err := t.fromRows(rows)
		if err != nil {
			return tenders, fmt.Errorf("%s:%w", op, err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (t *Repository) CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error) {
	const op = "repository.tender.CheckResponsible"
	tender := t.db.QueryRowxContext(ctx, `SELECT COUNT(*) FROM tender left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1 AND tender.id = $2 and version = (SELECT MAX(version) FROM tender t2 WHERE tender.id = t2.id)`, username, tenderId)

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

func (t *Repository) FetchById(ctx context.Context, tenderId string) (model.Tender, error) {
	const op = "repository.tender.FetchById"

	q := fmt.Sprintf(`SELECT %s FROM tender WHERE tender.id = $1 and version = (SELECT MAX(version) FROM tender t2 WHERE tender.id = t2.id)`, strings.Join(column, ","))
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

func (t *Repository) Create(ctx context.Context, saveModel model.Tender) (model.Tender, error) {
	const op = "repository.tender.Create"

	r := toRow(saveModel)

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
	saveModel.Id = id

	return saveModel, nil
}

func (t *Repository) Edite(ctx context.Context, tenderNew model.Tender, tender model.Tender) (model.Tender, error) {
	const op = "repository.tender.Edite"
	var r row
	tenderNew.Version = tender.Version + 1

	r = toRow(tenderNew)
	q := `INSERT INTO tender (id, name, description, service_type, status, version, responsible) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	result, err := t.db.QueryxContext(ctx, q, r.Id, r.Name, r.Description, r.ServiceType, r.Status, r.Version, r.Responsible)
	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	err = result.Close()
	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	return tenderNew, nil
}

func (t *Repository) Rollback(ctx context.Context, id string, version string) (model.Tender, error) {
	const op = "repository.tender.Rollback"

	q := `INSERT INTO tender (id, name, description, service_type, status,version,responsible) SELECT id, name, description, service_type, status, (SELECT version + 1 FROM tender WHERE tender.id = $1 and version = (SELECT MAX(version) FROM tender t2 WHERE tender.id = t2.id)) as version, responsible FROM tender WHERE id = $1 and version = $2 RETURNING id`
	tender := t.db.QueryRowxContext(ctx, q, id, version)
	err := tender.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return model.Tender{}, model.NotFound
	}
	if err != nil {
		return model.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	var tenderId string
	err = tender.Scan(&tenderId)
	if err != nil {
		return model.Tender{}, model.NotFound
	}

	resp, err := t.FetchById(ctx, tenderId)
	if err != nil {
		return model.Tender{}, model.NotFound
	}

	return resp, nil
}

var column = []string{"tender.id", "tender.name", "tender.description", "tender.service_type", "tender.status", "tender.version", "tender.responsible"}

func (t *Repository) fromRows(rows *sqlx.Rows) (model.Tender, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
