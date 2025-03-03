package bid_feedback

import (
	"TenderServiceApi/internal/model"
)

type row struct {
	Id          string `db:"id"`
	BidID       string `db:"bid_id"`
	Description string `db:"description"`
	Responsible string `db:"responsible"`
	CreatedAt   string `db:"created_at"`
}

func (r *row) toModel() model.BidFeedback {
	return model.BidFeedback{
		Id:          r.Id,
		BidID:       r.BidID,
		Description: r.Description,
		Responsible: r.Responsible,
		CreatedAt:   r.CreatedAt,
	}
}

func toRow(m model.BidFeedback) row {
	return row{
		Id:          m.Id,
		BidID:       m.BidID,
		Description: m.Description,
		Responsible: m.Responsible,
		CreatedAt:   m.CreatedAt,
	}
}
