package bid_decision

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

func (t *Repository) Save(ctx context.Context, bidID string, decision string, responsible string) error {
	const op = "repository.bids.SubmitDecision"

	q := `INSERT INTO decision_bid (bid_id, decision, responsible) VALUES ($1,$2,$3)`
	result, err := t.db.QueryxContext(ctx, q, bidID, decision, responsible)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = result.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (t *Repository) FetchByBidID(ctx context.Context, bidID string) ([]model.BidsDecisions, error) {
	const op = "repository.bids.CheckQuorum"
	var decisions []model.BidsDecisions

	q := fmt.Sprintf(`SELECT %s  FROM decision_bid WHERE decision_bid.bid_id = $1`, strings.Join(column, ","))

	rows, err := t.db.QueryxContext(ctx, q, bidID)
	if errors.Is(err, sql.ErrNoRows) {
		return decisions, model.NotFound
	}
	if err != nil {
		return decisions, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		bidsDecision, err := t.fromRows(rows)
		if err != nil {
			return decisions, fmt.Errorf("%s:%w", op, err)
		}
		decisions = append(decisions, bidsDecision)
	}

	return decisions, nil
}

var column = []string{"decision_bid.id", "decision_bid.bid_id", "decision_bid.decision", "decision_bid.responsible"}

func (t *Repository) fromRows(rows *sqlx.Rows) (model.BidsDecisions, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
