package bids

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

func (t *Repository) FetchListByTender(ctx context.Context, tenderId string) ([]model.Bids, error) {
	const op = "repository.bids.FetchList"
	var bids []model.Bids
	var rows *sqlx.Rows
	var err error

	q := fmt.Sprintf(`SELECT %s FROM bids left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id and bids.tender_id = $1) and bids.status = ANY ('{"CANCELED","PUBLISHED"}')`, strings.Join(column, ","))
	rows, err = t.db.QueryxContext(ctx, q, tenderId)
	if errors.Is(err, sql.ErrNoRows) {
		return bids, model.NotFound
	}
	if err != nil {
		return bids, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		tender, err := t.fromRows(rows)
		if err != nil {
			return bids, fmt.Errorf("%s:%w", op, err)
		}
		bids = append(bids, tender)
	}

	return bids, nil
}

func (t *Repository) FetchListByUser(ctx context.Context, username string) ([]model.Bids, error) {
	const op = "repository.bids.FetchListByUser"
	var bids []model.Bids
	var rows *sqlx.Rows
	var err error

	q := fmt.Sprintf(`SELECT %s FROM bids left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1 and version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id)`, strings.Join(column, ","))
	rows, err = t.db.QueryxContext(ctx, q, username)
	if errors.Is(err, sql.ErrNoRows) {
		return bids, model.NotFound
	}
	if err != nil {
		return bids, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		tender, err := t.fromRows(rows)
		if err != nil {
			return bids, fmt.Errorf("%s:%w", op, err)
		}
		bids = append(bids, tender)
	}

	return bids, nil
}

func (t *Repository) Create(ctx context.Context, saveModel model.Bids) (model.Bids, error) {
	const op = "repository.bids.Create"

	r := toRow(saveModel)

	q := "INSERT INTO bids (name, description, status, tender_id, version, responsible) VALUES($1,$2,$3,$4,$5,$6) RETURNING id"

	result := t.db.QueryRowxContext(ctx, q, r.Name, r.Description, r.Status, r.TenderId, r.Version, r.Responsible)

	err := result.Err()

	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	var id string

	err = result.Scan(&id)
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}
	saveModel.Id = id

	return saveModel, nil
}

func (t *Repository) CheckResponsible(ctx context.Context, username string, bidsId string) (bool, error) {
	const op = "repository.bids.CheckResponsible"
	tender := t.db.QueryRowxContext(ctx, `SELECT COUNT(*) FROM bids left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1 AND bids.id = $2 and version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id)`, username, bidsId)
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
	return false, model.NotFindResponsible
}

func (t *Repository) FetchById(ctx context.Context, bidsId string) (model.Bids, error) {
	const op = "repository.tender.FetchById"

	q := fmt.Sprintf(`SELECT %s FROM bids WHERE bids.id = $1 and version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id)`, strings.Join(column, ","))
	tender := t.db.QueryRowxContext(ctx, q, bidsId)
	err := tender.Err()

	if errors.Is(err, sql.ErrNoRows) {
		return model.Bids{}, model.NotFound
	}
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	var te row
	err = tender.StructScan(&te)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Bids{}, model.NotFound
	}
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	return te.toModel(), nil
}

func (t *Repository) UpdateStatus(ctx context.Context, bidID string, status string) (model.Bids, error) {
	const op = "repository.bids.UpdateStatus"

	q := `UPDATE bids SET status = $1 WHERE id = $2`
	result, err := t.db.QueryxContext(ctx, q, status, bidID)
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	err = result.Close()
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := t.FetchById(ctx, bidID)
	if err != nil {
		return model.Bids{}, model.NotFound
	}

	return resp, nil
}

func (t *Repository) Edit(ctx context.Context, bidNew model.Bids, bid model.Bids) (model.Bids, error) {
	const op = "repository.bids.Edit"
	var r row
	bidNew.Version = bid.Version + 1

	r = toRow(bidNew)
	q := `INSERT INTO bids (id, name, description, status, tender_id, version, responsible) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	result, err := t.db.QueryxContext(ctx, q, r.Id, r.Name, r.Description, r.Status, r.TenderId, r.Version, r.Responsible)
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	err = result.Close()
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	return bidNew, nil
}

func (t *Repository) Rollback(ctx context.Context, id string, version string) (model.Bids, error) {
	const op = "repository.bids.Rollback"

	q := `INSERT INTO bids (id, name, description, status, tender_id, version, responsible) SELECT id, name, description, status, tender_id, (SELECT version + 1 FROM bids WHERE bids.id = $1 and version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id)) as version, responsible FROM bids WHERE id = $1 and version = $2 RETURNING id`
	bid := t.db.QueryRowxContext(ctx, q, id, version)
	err := bid.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return model.Bids{}, fmt.Errorf("%s: %w", op, model.NotFound)
	}
	if err != nil {
		return model.Bids{}, fmt.Errorf("%s: %w", op, err)
	}

	var bidID string
	err = bid.Scan(&bidID)
	if err != nil {
		return model.Bids{}, model.NotFound
	}

	resp, err := t.FetchById(ctx, bidID)
	if err != nil {
		return model.Bids{}, model.NotFound
	}

	return resp, nil
}

var column = []string{"bids.id", "bids.name", "bids.description", "bids.status", "bids.tender_id", "bids.version", "bids.responsible"}

func (t *Repository) fromRows(rows *sqlx.Rows) (model.Bids, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
