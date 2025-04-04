package edit

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	bid        repository
	bidsVerify useCaseBidVerify
}

//go:generate mockery --inpackage --name=repository --exported --testonly --inpackage-suffix
type repository interface {
	Edit(ctx context.Context, bidNew model.Bids, bid model.Bids) (model.Bids, error)
	Rollback(ctx context.Context, id string, version string) (model.Bids, error)
	FetchById(ctx context.Context, bidID string) (model.Bids, error)
	UpdateStatus(ctx context.Context, bidID string, status string) (model.Bids, error)
}

//go:generate mockery --inpackage --name=useCaseBidVerify --exported --testonly --inpackage-suffix
type useCaseBidVerify interface {
	CheckResponsible(ctx context.Context, username string, bidID string) (bool, error)
}

func NewService(r repository, useCaseBidVerify useCaseBidVerify) *Service {
	return &Service{bid: r, bidsVerify: useCaseBidVerify}
}

func (s *Service) Edit(ctx context.Context, id string, username string, bidNew model.Bids) (model.Bids, error) {
	_, err := s.bidsVerify.CheckResponsible(ctx, username, id)
	if err != nil {
		return model.Bids{}, err
	}

	bid, err := s.bid.FetchById(ctx, id)
	if err != nil {
		return model.Bids{}, err
	}

	bidNew.FillDefault(bid)
	bidNew.Version = bid.Version + 1
	resp, err := s.bid.Edit(ctx, bidNew, bid)
	if err != nil {
		return model.Bids{}, err
	}

	return resp, nil
}

func (s *Service) Rollback(ctx context.Context, id string, username string, version string) (model.Bids, error) {
	bid, err := s.bid.FetchById(ctx, id)
	if err != nil {
		return model.Bids{}, err
	}

	_, err = s.bidsVerify.CheckResponsible(ctx, username, bid.Id)
	if err != nil {
		return model.Bids{}, err
	}

	resp, err := s.bid.Rollback(ctx, id, version)
	if err != nil {
		return model.Bids{}, err
	}

	return resp, nil
}

func (s *Service) Status(ctx context.Context, username string, bidID string, status string) (model.Bids, error) {
	_, err := s.bidsVerify.CheckResponsible(ctx, username, bidID)
	if err != nil {
		return model.Bids{}, err
	}

	bid, err := s.bid.UpdateStatus(ctx, bidID, status)
	if err != nil {
		return model.Bids{}, err
	}

	return bid, nil
}
