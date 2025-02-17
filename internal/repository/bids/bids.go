package tender

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
		fmt.Println(err) // nolint:all
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

var column = []string{"bids.id", "bids.name", "bids.description", "bids.status", "bids.tender_id", "bids.version", "bids.responsible"}

func (t *Repository) fromRows(rows *sqlx.Rows) (model.Bids, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
