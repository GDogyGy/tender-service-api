package convert

import (
	"TenderServiceApi/internal/handlers/types/transport"
	"TenderServiceApi/internal/model"
)

func BidsModelToTransport(m model.Bids) transport.Bid {
	return transport.Bid{
		Id:          m.Id,
		Description: m.Description,
		Name:        m.Name,
		Responsible: m.Responsible,
		Status:      transport.BidStatus(m.Status),
		TenderId:    m.TenderId,
		Version:     m.Version,
	}
}

func BidsTransportToModel(t transport.Bid) model.Bids {
	return model.Bids{
		Id:          t.Id,
		Description: t.Description,
		Name:        t.Name,
		Responsible: t.Responsible,
		Status:      string(t.Status),
		TenderId:    t.TenderId,
		Version:     t.Version,
	}
}
