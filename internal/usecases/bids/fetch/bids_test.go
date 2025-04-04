package fetch_test

import (
	"context"
	"fmt"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/bids/fetch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseFetchListByTender(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		tenderID     string
		prepare      func(repository *Mockrepository, verify *MockuseCaseTenderVerify)
		expectations func(t *testing.T, result []model.Bids, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repository.On("FetchListByTender", mock.Anything, "1").Return([]model.Bids{
					{
						Id:          "1",
						Name:        "test",
						Description: "test",
						Status:      "PUBLISHED",
						TenderId:    "1",
						Version:     1,
						Responsible: "1",
					},
				}, nil)
			},
			expectations: func(t *testing.T, result []model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result []model.Bids, err error) {
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
			verifyTender := NewMockuseCaseTenderVerify(t)
			verifyBids := NewMockuseCaseBidsVerify(t)
			bidFeedback := NewMockbidFeedback(t)

			tc.prepare(repository, verifyTender)

			useCase := fetch.NewService(repository, bidFeedback, verifyTender, verifyBids)
			r, err := useCase.FetchListByTender(tc.ctx, tc.username, tc.tenderID)

			tc.expectations(t, r, err)
		})
	}
}

func TestUseCaseFetchStatus(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		tenderID     string
		prepare      func(repository *Mockrepository, verify *MockuseCaseBidsVerify)
		expectations func(t *testing.T, result model.Bids, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidsVerify) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "test",
					Description: "test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidsVerify) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},

		{
			name:     "Error FetchById",
			ctx:      context.Background(),
			username: "user1",
			tenderID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidsVerify) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				repository.On("FetchById", mock.Anything, "1").Return(model.Bids{}, model.NotFound)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
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
			verifyTender := NewMockuseCaseTenderVerify(t)
			verifyBids := NewMockuseCaseBidsVerify(t)
			bidFeedback := NewMockbidFeedback(t)

			tc.prepare(repository, verifyBids)

			useCase := fetch.NewService(repository, bidFeedback, verifyTender, verifyBids)
			r, err := useCase.FetchStatus(tc.ctx, tc.username, tc.tenderID)

			tc.expectations(t, r, err)
		})
	}
}

func TestUseCaseFetchReviews(t *testing.T) {
	cases := []struct {
		name           string
		ctx            context.Context
		username       string
		authorUsername string
		tenderID       string
		organizationID string
		prepare        func(repository *Mockrepository, verify *MockuseCaseTenderVerify, bidFeedback *MockbidFeedback)
		expectations   func(t *testing.T, result []model.BidFeedback, err error)
	}{
		{
			name:           "Success",
			ctx:            context.Background(),
			username:       "user1",
			authorUsername: "user1",
			tenderID:       "1",
			organizationID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, bidFeedback *MockbidFeedback) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				bidFeedback.On("FetchReviews", mock.Anything, "1", "user1", "1").Return([]model.BidFeedback{{
					Id:          "1",
					BidID:       "1",
					Description: "test",
					Responsible: "1",
					CreatedAt:   "",
				},
				}, nil)
			},
			expectations: func(t *testing.T, result []model.BidFeedback, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
			},
		},
		{
			name:           "Error CheckResponsible",
			ctx:            context.Background(),
			username:       "user1",
			authorUsername: "user1",
			tenderID:       "1",
			organizationID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, bidFeedback *MockbidFeedback) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result []model.BidFeedback, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},

		{
			name:           "Error FetchReviews",
			ctx:            context.Background(),
			username:       "user1",
			authorUsername: "user2",
			tenderID:       "1",
			organizationID: "1",
			prepare: func(repository *Mockrepository, verify *MockuseCaseTenderVerify, bidFeedback *MockbidFeedback) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				bidFeedback.On("FetchReviews", mock.Anything, "1", "user2", "1").Return([]model.BidFeedback{}, fmt.Errorf(""))
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
			verifyTender := NewMockuseCaseTenderVerify(t)
			verifyBids := NewMockuseCaseBidsVerify(t)
			bidFeedback := NewMockbidFeedback(t)

			tc.prepare(repository, verifyTender, bidFeedback)

			useCase := fetch.NewService(repository, bidFeedback, verifyTender, verifyBids)
			r, err := useCase.FetchReviews(tc.ctx, tc.username, tc.tenderID, tc.authorUsername, tc.organizationID)

			tc.expectations(t, r, err)
		})
	}
}
