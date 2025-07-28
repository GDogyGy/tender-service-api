package create_test

import (
	"context"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/bid_feedback/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseCreate(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		bidID        string
		saveModel    model.BidFeedback
		prepare      func(repositoryFeedback *MockrepositoryFeedback, repositoryOrganization *MockrepositoryOrganization, verify *MockuseCaseBidFeedbackVerify, saveModel model.BidFeedback)
		expectations func(t *testing.T, result model.BidFeedback, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			bidID:    "1",
			saveModel: model.BidFeedback{
				BidID:       "1",
				Description: "Test",
				Responsible: "1",
			},
			prepare: func(repositoryFeedback *MockrepositoryFeedback, repositoryOrganization *MockrepositoryOrganization, verify *MockuseCaseBidFeedbackVerify, saveModel model.BidFeedback) {

				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repositoryOrganization.On("FetchByUserName", mock.Anything, "user1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				repositoryFeedback.On("Create", mock.Anything, saveModel).Return(model.BidFeedback{
					Id:          "1",
					BidID:       "1",
					Description: "Test",
					Responsible: "1",
					CreatedAt:   "10-10-2025",
				}, nil)
			},
			expectations: func(t *testing.T, result model.BidFeedback, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
			},
		},
		{
			name:     "Error useCaseBidFeedbackVerify CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			bidID:    "1",
			saveModel: model.BidFeedback{
				BidID:       "1",
				Description: "Test",
				Responsible: "1",
			},
			prepare: func(repositoryFeedback *MockrepositoryFeedback, repositoryOrganization *MockrepositoryOrganization, verify *MockuseCaseBidFeedbackVerify, saveModel model.BidFeedback) {

				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.BidFeedback, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error repositoryOrganization CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			bidID:    "1",
			saveModel: model.BidFeedback{
				BidID:       "1",
				Description: "Test",
				Responsible: "1",
			},
			prepare: func(repositoryFeedback *MockrepositoryFeedback, repositoryOrganization *MockrepositoryOrganization, verify *MockuseCaseBidFeedbackVerify, saveModel model.BidFeedback) {

				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repositoryOrganization.On("FetchByUserName", mock.Anything, "user1").Return(model.OrganizationResponsible{}, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.BidFeedback, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repositoryFeedback := NewMockrepositoryFeedback(t)
			repositoryOrganization := NewMockrepositoryOrganization(t)

			useCaseBidFeedbackVerify := NewMockuseCaseBidFeedbackVerify(t)

			tc.prepare(repositoryFeedback, repositoryOrganization, useCaseBidFeedbackVerify, tc.saveModel)

			useCase := create.NewService(repositoryFeedback, repositoryOrganization, useCaseBidFeedbackVerify)
			r, err := useCase.Create(tc.ctx, tc.username, tc.bidID, tc.saveModel)

			tc.expectations(t, r, err)
		})
	}
}
