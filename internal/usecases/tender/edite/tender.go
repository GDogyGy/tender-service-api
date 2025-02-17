package edite

import (
	"context"
	"reflect"

	"TenderServiceApi/internal/model"
)

type Service struct {
	tender              repository
	useCaseTenderVerify useCaseTenderVerify
}

type repository interface {
	Edite(ctx context.Context, tenderNew model.Tender, tender model.Tender) (model.Tender, error)
	FetchById(ctx context.Context, tenderId string) (model.Tender, error)
	Rollback(ctx context.Context, id string, version string) (model.Tender, error)
	UpdateStatus(ctx context.Context, tenderId string, status string) (model.Tender, error)
}

type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

func NewService(r repository, useCaseTenderVerify useCaseTenderVerify) *Service {
	return &Service{tender: r, useCaseTenderVerify: useCaseTenderVerify}
}

func (s *Service) copyNonEmptyValues(src, dst interface{}) {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// Убедимся, что переданные значения - это указатели на структуры
	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return
	}

	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()

	// Проходим по полям структуры
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.Field(i)

		// Проверяем, является ли поле пустым
		isEmpty := reflect.DeepEqual(dstField.Interface(), reflect.Zero(dstField.Type()).Interface())
		if isEmpty {
			// Копируем значение из src в dst, если dst пуст
			dstField.Set(srcField)
		}
	}
}

func (s *Service) Edite(ctx context.Context, id string, username string, tenderNew model.Tender) (model.Tender, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, id)
	if err != nil {
		return model.Tender{}, err
	}

	tender, err := s.tender.FetchById(ctx, id)
	if err != nil {
		return model.Tender{}, err
	}

	s.copyNonEmptyValues(&tender, &tenderNew) // TODO: Спросить у димы более изящный способ

	resp, err := s.tender.Edite(ctx, tenderNew, tender)
	if err != nil {
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
