package edit

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	tender              repository
	useCaseTenderVerify useCaseTenderVerify
}

//go:generate mockery --inpackage --name=repository --exported --testonly --inpackage-suffix
type repository interface {
	Edit(ctx context.Context, tenderNew model.Tender, tender model.Tender) (model.Tender, error)
	FetchById(ctx context.Context, tenderId string) (model.Tender, error)
	Rollback(ctx context.Context, id string, version string) (model.Tender, error)
	UpdateStatus(ctx context.Context, tenderId string, status string) (model.Tender, error)
}

//go:generate mockery --inpackage --name=useCaseTenderVerify --exported --testonly --inpackage-suffix
type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

func NewService(r repository, useCaseTenderVerify useCaseTenderVerify) *Service {
	return &Service{tender: r, useCaseTenderVerify: useCaseTenderVerify}
}

func (s *Service) Edit(ctx context.Context, id string, username string, tenderNew model.Tender) (model.Tender, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, id)
	if err != nil {
		return model.Tender{}, err
	}

	tender, err := s.tender.FetchById(ctx, id)
	if err != nil {
		return model.Tender{}, err
	}

	tenderNew.FillDefault(tender)
	tenderNew.Version = tender.Version + 1

	resp, err := s.tender.Edit(ctx, tenderNew, tender)
	if err != nil { // TODO: оказывается эту штуку не так просто проверить
		return model.Tender{}, err
	}

	return resp, nil
}

func (s *Service) Rollback(ctx context.Context, id string, username string, version string) (model.Tender, error) {
	tender, err := s.tender.FetchById(ctx, id)
	if err != nil {
		return model.Tender{}, err
	}

	_, err = s.useCaseTenderVerify.CheckResponsible(ctx, username, tender.Id)
	if err != nil {
		return model.Tender{}, err
	}

	resp, err := s.tender.Rollback(ctx, id, version)
	if err != nil {
		return model.Tender{}, err
	}

	return resp, nil
}

func (s *Service) Status(ctx context.Context, username string, tenderId string, status string) (model.Tender, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, tenderId)
	if err != nil {
		return model.Tender{}, err
	}

	tender, err := s.tender.UpdateStatus(ctx, tenderId, status)
	if err != nil {
		return model.Tender{}, err
	}

	return tender, nil
}
