package edite

import (
	"context"
	"reflect"

	"TenderServiceApi/internal/model"
)

type Service struct {
	tender                    repository
	useCaseOrganizationVerify useCaseOrganizationVerify
	useCaseOrganizationFetch  useCaseOrganizationFetch
}

type repository interface {
	Edite(ctx context.Context, tenderNew model.Tender, tender model.Tender) (model.Tender, error)
	FetchById(ctx context.Context, tenderId string) (model.Tender, error)
	Rollback(ctx context.Context, id string, version string) (model.Tender, error)
}

type useCaseOrganizationVerify interface {
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
}

type useCaseOrganizationFetch interface {
	FetchRelationsById(ctx context.Context, id string) (model.OrganizationResponsible, error)
}

func NewService(r repository, useCaseOrganizationVerify useCaseOrganizationVerify, useCaseOrganizationFetch useCaseOrganizationFetch) *Service {
	return &Service{tender: r, useCaseOrganizationVerify: useCaseOrganizationVerify, useCaseOrganizationFetch: useCaseOrganizationFetch}
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
	tender, err := s.tender.FetchById(ctx, id)
	if err != nil {
		return model.Tender{}, err
	}

	organization, err := s.useCaseOrganizationFetch.FetchRelationsById(ctx, tender.Responsible)
	if err != nil {
		return model.Tender{}, err
	}

	s.copyNonEmptyValues(&tender, &tenderNew) // TODO: Спросить у димы более изящный способ

	_, err = s.useCaseOrganizationVerify.CheckResponsible(ctx, username, organization.OrganizationId)
	if err != nil {
		return model.Tender{}, err
	}

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

	// TODO: Можно избавиться от запроса но тогда придется прокидывать useCaseTenderVerification и зависимостей больше 3 тут получится
	organization, err := s.useCaseOrganizationFetch.FetchRelationsById(ctx, tender.Responsible)
	if err != nil {
		return model.Tender{}, err
	}

	_, err = s.useCaseOrganizationVerify.CheckResponsible(ctx, username, organization.OrganizationId)
	if err != nil {
		return model.Tender{}, err
	}

	resp, err := s.tender.Rollback(ctx, id, version)
	if err != nil {
		return model.Tender{}, err
	}

	return resp, nil
}
