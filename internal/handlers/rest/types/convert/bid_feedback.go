package convert

import (
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"TenderServiceApi/internal/model"
)

func BidFeedbackModelToTransport(m model.BidFeedback) transport.Feedback {
	return transport.Feedback{
		Id:          m.Id,
		BidId:       m.BidID,
		Description: m.Description,
		Responsible: m.Responsible,
		CreatedAt:   m.CreatedAt,
	}
}

func BidFeedbackTransportToModel(f transport.Feedback) model.BidFeedback {
	return model.BidFeedback{
		Id:          f.Id,
		BidID:       f.BidId,
		Description: f.Description,
		Responsible: f.Responsible,
		CreatedAt:   f.CreatedAt,
	}
}
