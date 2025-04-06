package edit_test

import (
	"context"
	"fmt"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/tender/edit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseEdit(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		id           string
		tender       model.Tender
		tenderNew    model.Tender
		prepare      func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tenderNew model.Tender, tender model.Tender)
		expectations func(t *testing.T, result model.Tender, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			id:       "1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     1,
				Responsible: "1",
			},
			tenderNew: model.Tender{
				Name:        "New Name",
				Description: "New Description",
				ServiceType: "New ServiceType",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tenderNew model.Tender, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)
				repository.On("FetchById", mock.Anything, "1").Return(tender, nil)
				_ = tenderNew
				tenderNew = model.Tender{
					Id:          "1",
					Name:        "New Name",
					Description: "New Description",
					ServiceType: "New ServiceType",
					Status:      "PUBLISHED",
					Version:     2,
					Responsible: "1",
				}

				repository.On("Edit", mock.Anything, tenderNew, tender).Return(tenderNew, nil)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, 2, result.Version)
			},
		},
		{
			name:      "Error CheckResponsible",
			ctx:       context.Background(),
			username:  "user1",
			id:        "1",
			tender:    model.Tender{},
			tenderNew: model.Tender{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tenderNew model.Tender, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, fmt.Errorf("error CheckResponsible"))
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Error(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name:      "Error FetchById",
			ctx:       context.Background(),
			username:  "user1",
			id:        "1",
			tender:    model.Tender{},
			tenderNew: model.Tender{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tenderNew model.Tender, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)
				repository.On("FetchById", mock.Anything, "1").Return(model.Tender{}, model.NotFound)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Error(t, err)
				assert.Empty(t, result)
			},
		},

		{
			name:     "Error Edit",
			ctx:      context.Background(),
			username: "user1",
			id:       "1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     1,
				Responsible: "1",
			},
			tenderNew: model.Tender{
				Name:        "New Name",
				Description: "New Description",
				ServiceType: "New ServiceType",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tenderNew model.Tender, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)
				repository.On("FetchById", mock.Anything, "1").Return(tender, nil)
				tenderNew.FillDefault(tender)
				tenderNew.Version = tender.Version + 1
				repository.On("Edit", mock.Anything, tenderNew, tender).Return(model.Tender{}, fmt.Errorf("error edit"))
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Error(t, err)
				assert.Empty(t, result)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repository := NewMockrepository(t)
			useCaseTenderVerify := NewMockuseCaseTenderVerify(t)

			tc.prepare(repository, useCaseTenderVerify, tc.tenderNew, tc.tender)

			useCase := edit.NewService(repository, useCaseTenderVerify)
			r, err := useCase.Edit(tc.ctx, tc.id, tc.username, tc.tenderNew)

			tc.expectations(t, r, err)
		})
	}
}

func TestUseCaseRollback(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		tender       model.Tender
		version      string
		prepare      func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender)
		expectations func(t *testing.T, result model.Tender, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			version: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				repository.On("FetchById", mock.Anything, "1").Return(tender, nil)
				verify.On("CheckResponsible", mock.Anything, "user1", tender.Id).Return(true, nil)
				var r model.Tender
				r.FillDefault(tender)
				r.Version = 1
				repository.On("Rollback", mock.Anything, tender.Id, "1").Return(r, nil)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Version)
			},
		},
		{
			name:     "Error FetchById",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			version: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				repository.On("FetchById", mock.Anything, "1").Return(model.Tender{}, model.NotFound)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			version: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				repository.On("FetchById", mock.Anything, "1").Return(tender, nil)
				verify.On("CheckResponsible", mock.Anything, "user1", tender.Id).Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error Rollback",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			version: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				repository.On("FetchById", mock.Anything, "1").Return(tender, nil)
				verify.On("CheckResponsible", mock.Anything, "user1", tender.Id).Return(true, nil)
				var r model.Tender
				r.FillDefault(tender)
				r.Version = 1
				repository.On("Rollback", mock.Anything, tender.Id, "1").Return(model.Tender{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repository := NewMockrepository(t)
			useCaseTenderVerify := NewMockuseCaseTenderVerify(t)

			tc.prepare(repository, useCaseTenderVerify, tc.tender)

			useCase := edit.NewService(repository, useCaseTenderVerify)
			r, err := useCase.Rollback(tc.ctx, tc.tender.Id, tc.username, tc.version)

			tc.expectations(t, r, err)
		})
	}
}

func TestUseCaseStatus(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		tender       model.Tender
		status       string
		prepare      func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender, status string)
		expectations func(t *testing.T, result model.Tender, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			status: "CREATED",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender, status string) {
				verify.On("CheckResponsible", mock.Anything, "user1", tender.Id).Return(true, nil)
				tender.Status = status
				repository.On("UpdateStatus", mock.Anything, tender.Id, status).Return(tender, nil)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, "CREATED", result.Status)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			status: "CREATED",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender, status string) {
				verify.On("CheckResponsible", mock.Anything, "user1", tender.Id).Return(false, model.NotFound)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error UpdateStatus",
			ctx:      context.Background(),
			username: "user1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "PUBLISHED",
				Version:     2,
				Responsible: "1",
			},
			status: "CREATED",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender, status string) {
				verify.On("CheckResponsible", mock.Anything, "user1", tender.Id).Return(true, nil)
				tender.Status = status
				repository.On("UpdateStatus", mock.Anything, tender.Id, status).Return(model.Tender{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repository := NewMockrepository(t)
			useCaseTenderVerify := NewMockuseCaseTenderVerify(t)

			tc.prepare(repository, useCaseTenderVerify, tc.tender, tc.status)

			useCase := edit.NewService(repository, useCaseTenderVerify)
			r, err := useCase.Status(tc.ctx, tc.username, tc.tender.Id, tc.status)

			tc.expectations(t, r, err)
		})
	}
}
