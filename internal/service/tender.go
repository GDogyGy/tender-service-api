package service

import "TenderServiceApi/internal/model"

type TenderService struct {
	tenderRepository TenderRepository
}

type TenderRepository interface {
	GetTenderList(serviceType string) ([]model.Tender, error)
}

func NewTenderService(tenderRepository TenderRepository) *TenderService {
	return &TenderService{tenderRepository: tenderRepository}
}

func (s *TenderService) GetTenderList(serviceType string) ([]model.Tender, error) {
	return s.tenderRepository.GetTenderList(serviceType)
}
