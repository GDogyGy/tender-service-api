package fetch_test

import (
	"context"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/tender/fetch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseFetchStatus(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		tenderID     string
		tender       model.Tender
		prepare      func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tenderNew model.Tender)
		expectations func(t *testing.T, result model.Tender, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			tender: model.Tender{
				Id:          "1",
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "CREATED",
				Version:     1,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repository.On("FetchById", mock.Anything, "1").Return(tender, nil)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			tender:   model.Tender{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
				assert.Error(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name:     "Error Fetch Tender",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			tender:   model.Tender{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, tender model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repository.On("FetchById", mock.Anything, "1").Return(model.Tender{}, model.NotFound)
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
			verify := NewMockuseCaseTenderVerify(t)

			tc.prepare(repository, verify, tc.tender)

			useCase := fetch.NewService(repository, verify)
			r, err := useCase.FetchStatus(tc.ctx, tc.username, tc.tenderID)

			tc.expectations(t, r, err)
		})
	}
}
