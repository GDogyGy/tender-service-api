package fetch_test

import (
	"context"
	"fmt"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/bid_feedback/fetch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseFetchReviews(t *testing.T) {
	cases := []struct {
		name           string
		ctx            context.Context
		username       string
		tenderID       string
		authorUsername string
		organizationID string
		prepare        func(repository *Mockrepository, tenderVerify *MockuseCaseTenderVerify, organizationVerify *MockuseCaseOrganizationVerify)
		expectations   func(t *testing.T, result []model.BidFeedback, err error)
	}{
		{
			name:           "Success",
			ctx:            context.Background(),
			username:       "user1",
			tenderID:       "1",
			authorUsername: "user2",
			organizationID: "1",
			prepare: func(repository *Mockrepository, tenderVerify *MockuseCaseTenderVerify, organizationVerify *MockuseCaseOrganizationVerify) {

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				repository.On("FetchReviews", mock.Anything, "1", "user2", "1").Return([]model.BidFeedback{
					{
						Id:          "1",
						BidID:       "1",
						Description: "Test",
						Responsible: "1",
						CreatedAt:   "10-10-2025",
					},
				}, nil)
			},
			expectations: func(t *testing.T, result []model.BidFeedback, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
			},
		},
		{
			name:           "Error useCaseTenderVerify CheckResponsible",
			ctx:            context.Background(),
			username:       "user1",
			tenderID:       "1",
			authorUsername: "user2",
			organizationID: "1",
			prepare: func(repository *Mockrepository, tenderVerify *MockuseCaseTenderVerify, organizationVerify *MockuseCaseOrganizationVerify) {

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result []model.BidFeedback, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error useCaseOrganizationVerify CheckResponsible",
			ctx:            context.Background(),
			username:       "user1",
			tenderID:       "1",
			authorUsername: "user2",
			organizationID: "1",
			prepare: func(repository *Mockrepository, tenderVerify *MockuseCaseTenderVerify, organizationVerify *MockuseCaseOrganizationVerify) {

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{}, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result []model.BidFeedback, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error feedbackRepo FetchReviews",
			ctx:            context.Background(),
			username:       "user1",
			tenderID:       "1",
			authorUsername: "user2",
			organizationID: "1",
			prepare: func(repository *Mockrepository, tenderVerify *MockuseCaseTenderVerify, organizationVerify *MockuseCaseOrganizationVerify) {

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				repository.On("FetchReviews", mock.Anything, "1", "user2", "1").Return([]model.BidFeedback{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result []model.BidFeedback, err error) {
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
			useCaseOrganizationVerify := NewMockuseCaseOrganizationVerify(t)

			tc.prepare(repository, useCaseTenderVerify, useCaseOrganizationVerify)

			useCase := fetch.NewService(repository, useCaseTenderVerify, useCaseOrganizationVerify)
			r, err := useCase.FetchReviews(tc.ctx, tc.username, tc.tenderID, tc.authorUsername, tc.organizationID)

			tc.expectations(t, r, err)
		})
	}
}
