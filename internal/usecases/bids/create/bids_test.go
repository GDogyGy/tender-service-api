package create_test

import (
	"context"
	"fmt"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/bids/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseCreate(t *testing.T) {
	cases := []struct {
		name           string
		ctx            context.Context
		username       string
		organizationID string
		bidsNew        model.Bids
		prepare        func(repository *Mockrepository, verify *MockuseCaseOrganizationVerify, bidsNew model.Bids)
		expectations   func(t *testing.T, result model.Bids, err error)
	}{
		{
			name:           "Success",
			ctx:            context.Background(),
			username:       "user1",
			organizationID: "123",
			bidsNew: model.Bids{
				Name:        "test",
				Description: "test",
				TenderId:    "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseOrganizationVerify, bidsNew model.Bids) {
				verify.On("CheckResponsible", mock.Anything, "user1", "123").Return(model.OrganizationResponsible{
					Id:             "123",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidsNew.Responsible = "123"
				repository.On("Create", mock.Anything, bidsNew).Return(model.Bids{
					Id:          "1",
					Name:        "test",
					Description: "test",
					TenderId:    "1",
					Version:     1,
					Responsible: "123",
				}, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Version)
			},
		},
		{
			name:           "Error CheckResponsible",
			ctx:            context.Background(),
			username:       "user1",
			organizationID: "1",
			bidsNew:        model.Bids{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseOrganizationVerify, bidsNew model.Bids) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{}, fmt.Errorf("error CheckResponsible"))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
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
			useCaseOrganizationVerify := NewMockuseCaseOrganizationVerify(t)

			tc.prepare(repository, useCaseOrganizationVerify, tc.bidsNew)

			useCase := create.NewService(repository, useCaseOrganizationVerify)
			r, err := useCase.Create(tc.ctx, tc.username, tc.organizationID, tc.bidsNew)

			tc.expectations(t, r, err)
		})
	}
}
