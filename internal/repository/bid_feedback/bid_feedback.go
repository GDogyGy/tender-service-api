package bid_feedback

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"TenderServiceApi/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (t *Repository) Create(ctx context.Context, saveModel model.BidFeedback) (model.BidFeedback, error) {
	const op = "repository.bidFeedback.Create"
	r := toRow(saveModel)
	currDate := time.Now()
	saveModel.CreatedAt = currDate.String()

	q := "INSERT INTO bid_feedback (bid_id, description, responsible, created_at) VALUES($1,$2,$3,$4) RETURNING id"

	result := t.db.QueryRowxContext(ctx, q, r.BidID, r.Description, r.Responsible, currDate)

	err := result.Err()
	if err != nil {
		return model.BidFeedback{}, fmt.Errorf("%s: %w", op, err)
	}

	var id string

	err = result.Scan(&id)
	if err != nil {
		return model.BidFeedback{}, fmt.Errorf("%s: %w", op, err)
	}

	saveModel.Id = id

	return saveModel, nil
}

func (t *Repository) FetchReviews(ctx context.Context, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error) {
	const op = "repository.bidFeedback.FetchReviews"
	var bidFeedbacks []model.BidFeedback
	var rows *sqlx.Rows
	var err error

	q := fmt.Sprintf(`SELECT %s FROM bids LEFT JOIN bid_feedback ON bid_feedback.bid_id = bids.id WHERE bids.tender_id = $1 and bids.version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id) and bids.responsible = ANY (SELECT organization_responsible.id FROM organization_responsible LEFT JOIN employee on employee.id = organization_responsible.user_id WHERE employee.username = $2 and organization_responsible.organization_id = $3) and bid_feedback.id is not null`, strings.Join(column, ","))
	rows, err = t.db.QueryxContext(ctx, q, tenderID, authorUsername, organizationID)
	if errors.Is(err, sql.ErrNoRows) {
		return bidFeedbacks, model.NotFound
	}
	if err != nil {
		return bidFeedbacks, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		bidFeedback, err := t.fromRows(rows)
		if err != nil {
			return bidFeedbacks, fmt.Errorf("%s:%w", op, err)
		}
		bidFeedbacks = append(bidFeedbacks, bidFeedback)
	}

	return bidFeedbacks, nil
}

func (t *Repository) CheckResponsible(ctx context.Context, username string, bidID string) (bool, error) {
	const op = "repository.bidFeedback.CheckResponsible"
	// TODO: Обсудить с димой: Улучшил запрос по сравнению с другими CheckResponsible в тендере
	bidFeedback := t.db.QueryRowxContext(ctx, `SELECT exists(SELECT * FROM bids
                         left join tender t on t.id = bids.tender_id
                         left join organization_responsible o on t.responsible = o.organization_id
                         left join employee e on o.user_id = e.id
 	WHERE e.username = $1
  	AND bids.id = $2 and bids.version = (SELECT MAX(version) FROM bids b2 WHERE bids.id = b2.id))`, username, bidID)

	err := bidFeedback.Err()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	c := false

	err = bidFeedback.Scan(&c)
	if errors.Is(err, sql.ErrNoRows) {
		return false, model.NotFound
	}
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if c {
		return true, nil
	}
	return false, model.NotFindResponsible
}

var column = []string{"bid_feedback.id", "bid_feedback.bid_id", "bid_feedback.description", "bid_feedback.responsible", "bid_feedback.created_at"}

func (t *Repository) fromRows(rows *sqlx.Rows) (model.BidFeedback, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
