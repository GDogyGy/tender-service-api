package tender

import (
	"TenderServiceApi/internal/model"
)

type row struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      string `db:"status"`
	TenderId    string `db:"tender_id"`
	Version     int    `db:"version"`
	Responsible string `db:"responsible"`
}

func (r *row) toModel() model.Bids {
	return model.Bids{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		TenderId:    r.TenderId,
		Version:     r.Version,
		Responsible: r.Responsible,
	}
}

func toRow(m model.Bids) row {
	return row{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Status:      m.Status,
		TenderId:    m.TenderId,
		Version:     m.Version,
		Responsible: m.Responsible,
	}
}
