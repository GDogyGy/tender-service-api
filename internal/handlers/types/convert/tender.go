package convert

import (
	"TenderServiceApi/internal/handlers/types/transport"
	"TenderServiceApi/internal/model"
)

func TenderModelToTransport(m model.Tender) transport.Tender {
	return transport.Tender{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		ServiceType: m.ServiceType,
		Status:      transport.TenderStatus(m.Status),
		Version:     m.Version,
		Responsible: m.Responsible,
	}
}

func TenderTransportToModel(t transport.Tender) model.Tender {
	return model.Tender{
		Id:          t.Id,
		Name:        t.Name,
		Description: t.Description,
		ServiceType: t.ServiceType,
		Status:      string(t.Status),
		Version:     t.Version,
		Responsible: t.Responsible,
	}
}
