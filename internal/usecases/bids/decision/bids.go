package decision

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	bid                bidRepository
	decision           decisionRepository
	tender             tenderRepository
	organizationVerify useCaseOrganizationVerify
	tenderVerify       useCaseTenderVerify
}

//go:generate mockery --inpackage --name=bidRepository --exported --testonly --inpackage-suffix
type bidRepository interface {
	FetchById(ctx context.Context, bidID string) (model.Bids, error)
	UpdateStatus(ctx context.Context, bidID string, status string) (model.Bids, error)
}

//go:generate mockery --inpackage --name=decisionRepository --exported --testonly --inpackage-suffix
type decisionRepository interface {
	Save(ctx context.Context, bidID string, decision string, responsible string) error
	FetchByBidID(ctx context.Context, bidID string) ([]model.BidsDecisions, error)
}

//go:generate mockery --inpackage --name=tenderRepository --exported --testonly --inpackage-suffix
type tenderRepository interface {
	UpdateStatus(ctx context.Context, tenderId string, status string) (model.Tender, error)
}

//go:generate mockery --inpackage --name=useCaseTenderVerify --exported --testonly --inpackage-suffix
type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

//go:generate mockery --inpackage --name=useCaseOrganizationVerify --exported --testonly --inpackage-suffix
type useCaseOrganizationVerify interface {
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
}

func NewService(r bidRepository, d decisionRepository, t tenderRepository, organizationVerify useCaseOrganizationVerify, useCaseTenderVerify useCaseTenderVerify) *Service {
	return &Service{bid: r, decision: d, tender: t, organizationVerify: organizationVerify, tenderVerify: useCaseTenderVerify}
}

func (s *Service) SubmitDecision(ctx context.Context, username string, bidID string, decision string, organizationID string) (model.Bids, error) {
	var bid model.Bids
	var quorume int

	// провека что такой ответственный вообще существует
	organization, err := s.organizationVerify.CheckResponsible(ctx, username, organizationID)
	if err != nil {
		return bid, err
	}

	bid, err = s.bid.FetchById(ctx, bidID)
	if err != nil {
		return bid, err
	}

	// провека что ответственный имеет право решать в этом тендере
	_, err = s.tenderVerify.CheckResponsible(ctx, username, bid.TenderId)
	if err != nil {
		return bid, err
	}

	// Достаем все решения
	decisions, err := s.decision.FetchByBidID(ctx, bidID)
	if err != nil {
		return bid, err
	}

	// проверяем нет ли решения от этого user ранее
	for i := range decisions {
		if decisions[i].Responsible == organization.Id {
			return bid, model.AlreadyVotedResponsible
		}

		if decisions[i].Decision == "APPROVED" {
			quorume++
		}
	}

	// Сохранение решения ответственного
	err = s.decision.Save(ctx, bidID, decision, organization.Id)
	if err != nil {
		return bid, err
	}

	// Проверяем решение = Reject и отклоняем предложение
	if decision == "REJECTED" {
		bid, err = s.bid.UpdateStatus(ctx, bidID, "CANCELED")
		if err != nil {
			return bid, err
		}
		return bid, nil
	}

	// Если есть согласование и кворум >= 3 то тендер закрывается
	if decision == "APPROVED" && quorume >= 2 {
		bid, err = s.bid.UpdateStatus(ctx, bidID, "APPROVED")
		if err != nil {
			return bid, err
		}

		_, err = s.tender.UpdateStatus(ctx, bid.TenderId, "CLOSED")
		if err != nil {
			return bid, err
		}
	}

	return bid, nil
}
