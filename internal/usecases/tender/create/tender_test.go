package create_test

import (
	"context"
	"fmt"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/tender/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseCreate(t *testing.T) {
	cases := []struct {
		name           string
		ctx            context.Context
		username       string
		organizationID string
		tenderNew      model.Tender
		prepare        func(repository *Mockrepository, verify *MockuseCaseOrganizationVerify, tenderNew model.Tender)
		expectations   func(t *testing.T, result model.Tender, err error)
	}{
		{
			name:           "Success",
			ctx:            context.Background(),
			username:       "user1",
			organizationID: "1",
			tenderNew: model.Tender{
				Name:        "test",
				Description: "test",
				ServiceType: "test",
				Status:      "CREATED",
				Version:     0,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseOrganizationVerify, tenderNew model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "123",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				tenderNew.Version = 1
				repository.On("Create", mock.Anything, tenderNew).Return(tenderNew, nil)
			},
			expectations: func(t *testing.T, result model.Tender, err error) {
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
			tenderNew:      model.Tender{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseOrganizationVerify, tenderNew model.Tender) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{}, fmt.Errorf("error CheckResponsible"))
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
			useCaseOrganizationVerify := NewMockuseCaseOrganizationVerify(t)

			tc.prepare(repository, useCaseOrganizationVerify, tc.tenderNew)

			useCase := create.NewService(repository, useCaseOrganizationVerify)
			r, err := useCase.Create(tc.ctx, tc.username, tc.organizationID, tc.tenderNew)

			tc.expectations(t, r, err)
		})
	}
}
