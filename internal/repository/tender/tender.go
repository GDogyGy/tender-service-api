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
	const op = "usecase.tender.GetTenderList"
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
	// TODO: defer rows.Close говорят в доке не нужен https://jmoiron.github.io/sqlx/ "The connection used by the Query remains..exhausted by the iteration via Next, or rows.Close() is called,"
	// TODO: тем более что под defer загнать нельзя так как ошибку возвращает"
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

var column = []string{"id", "name", "description", "service_type", "status", "responsible"}

func (t *Repository) tenderFromRows(rows *sqlx.Rows) (model.Tender, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
