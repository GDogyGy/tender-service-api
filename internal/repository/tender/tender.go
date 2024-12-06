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

func (t *Repository) FetchListByUser(ctx context.Context, username string) ([]model.Tender, error) {
	const op = "repository.tender.FetchListByUser"
	var tenders []model.Tender
	var rows *sqlx.Rows
	var err error

	// TODO: подумал что тут через facade не нужно делать, так как нет склейки моделей, простой join но на выходе ожидается model.Tender
	q := fmt.Sprintf(`SELECT %s FROM tender left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1`, strings.Join(column, ","))
	if err != nil {
		return tenders, model.NotFound
	}
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

var column = []string{"tender.id", "tender.name", "tender.description", "tender.service_type", "tender.status", "tender.responsible"}

func (t *Repository) tenderFromRows(rows *sqlx.Rows) (model.Tender, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
