package bid_decision

import (
	"TenderServiceApi/internal/model"
)

type row struct {
	Id          string `db:"id"`
	BidID       string `db:"bid_id"`
	Decision    string `db:"decision"`
	Responsible string `db:"responsible"`
}

func (r *row) toModel() model.BidsDecisions {
	return model.BidsDecisions{
		Id:          r.Id,
		BidID:       r.BidID,
		Decision:    r.Decision,
		Responsible: r.Responsible,
	}
}
