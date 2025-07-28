package edit_test

import (
	"context"
	"fmt"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/bids/edit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseEdit(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		id           string
		bid          model.Bids
		bidNew       model.Bids
		prepare      func(repository *Mockrepository, verify *MockuseCaseBidVerify, bidNew model.Bids, bid model.Bids)
		expectations func(t *testing.T, result model.Bids, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			id:       "1",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "PUBLISHED",
				TenderId:    "1",
				Version:     1,
				Responsible: "1",
			},
			bidNew: model.Bids{
				Name:        "New Name",
				Description: "New Description",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bidNew model.Bids, bid model.Bids) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)
				repository.On("FetchById", mock.Anything, "1").Return(bid, nil)

				bidNew.FillDefault(bid)
				bidNew.Version = 2
				repository.On("Edit", mock.Anything, bidNew, bid).Return(model.Bids{
					Id:          "1",
					Name:        "New Name",
					Description: "New Description",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     2,
					Responsible: "1",
				}, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, 2, result.Version)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			id:       "1",
			bid:      model.Bids{},
			bidNew:   model.Bids{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bidNew model.Bids, bid model.Bids) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(false, fmt.Errorf("error CheckResponsible"))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Error(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name:     "Error FetchById",
			ctx:      context.Background(),
			username: "user1",
			id:       "1",
			bid:      model.Bids{},
			bidNew:   model.Bids{},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bidNew model.Bids, bid model.Bids) {
				verify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)
				repository.On("FetchById", mock.Anything, "1").Return(model.Bids{}, model.NotFound)
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
			useCaseBidVerify := NewMockuseCaseBidVerify(t)

			tc.prepare(repository, useCaseBidVerify, tc.bidNew, tc.bid)

			useCase := edit.NewService(repository, useCaseBidVerify)
			r, err := useCase.Edit(tc.ctx, tc.id, tc.username, tc.bidNew)

			tc.expectations(t, r, err)
		})
	}
}

func TestUseCaseRollback(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		version      string
		bid          model.Bids
		prepare      func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids)
		expectations func(t *testing.T, result model.Bids, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			version:  "1",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "PUBLISHED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids) {
				repository.On("FetchById", mock.Anything, bid.Id).Return(bid, nil)

				verify.On("CheckResponsible", mock.Anything, "user1", bid.Id).Return(true, nil)

				repository.On("Rollback", mock.Anything, bid.Id, "1").Return(model.Bids{
					Id:          "1",
					Name:        "New Name",
					Description: "New Description",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Version)
			},
		},
		{
			name:     "Error FetchById",
			ctx:      context.Background(),
			username: "user1",
			version:  "1",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "PUBLISHED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids) {
				repository.On("FetchById", mock.Anything, bid.Id).Return(model.Bids{}, model.NotFound)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			version:  "1",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "PUBLISHED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids) {
				repository.On("FetchById", mock.Anything, bid.Id).Return(bid, nil)

				verify.On("CheckResponsible", mock.Anything, "user1", bid.Id).Return(false, model.NotFindResponsible)

			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error Rollback",
			ctx:      context.Background(),
			username: "user1",
			version:  "1",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "PUBLISHED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids) {
				repository.On("FetchById", mock.Anything, bid.Id).Return(bid, nil)

				verify.On("CheckResponsible", mock.Anything, "user1", bid.Id).Return(true, nil)

				repository.On("Rollback", mock.Anything, bid.Id, "1").Return(model.Bids{}, fmt.Errorf(""))
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
			useCaseBidVerify := NewMockuseCaseBidVerify(t)

			tc.prepare(repository, useCaseBidVerify, tc.bid)

			useCase := edit.NewService(repository, useCaseBidVerify)
			r, err := useCase.Rollback(tc.ctx, tc.bid.Id, tc.username, tc.version)

			tc.expectations(t, r, err)
		})
	}
}

func TestUseCaseStatus(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		username     string
		status       string
		bid          model.Bids
		prepare      func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids, status string)
		expectations func(t *testing.T, result model.Bids, err error)
	}{
		{
			name:     "Success",
			ctx:      context.Background(),
			username: "user1",
			status:   "PUBLISHED",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "CREATED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids, status string) {
				verify.On("CheckResponsible", mock.Anything, "user1", bid.Id).Return(true, nil)

				bid.Status = status
				repository.On("UpdateStatus", mock.Anything, bid.Id, status).Return(bid, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, "PUBLISHED", result.Status)
			},
		},
		{
			name:     "Error CheckResponsible",
			ctx:      context.Background(),
			username: "user1",
			status:   "PUBLISHED",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "CREATED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids, status string) {
				verify.On("CheckResponsible", mock.Anything, "user1", bid.Id).Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:     "Error UpdateStatus",
			ctx:      context.Background(),
			username: "user1",
			status:   "PUBLISHED",
			bid: model.Bids{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Status:      "CREATED",
				TenderId:    "1",
				Version:     2,
				Responsible: "1",
			},
			prepare: func(repository *Mockrepository, verify *MockuseCaseBidVerify, bid model.Bids, status string) {
				verify.On("CheckResponsible", mock.Anything, "user1", bid.Id).Return(true, nil)
				bid.Status = status
				repository.On("UpdateStatus", mock.Anything, bid.Id, status).Return(model.Bids{}, fmt.Errorf(""))
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
			useCaseBidVerify := NewMockuseCaseBidVerify(t)

			tc.prepare(repository, useCaseBidVerify, tc.bid, tc.status)

			useCase := edit.NewService(repository, useCaseBidVerify)
			r, err := useCase.Status(tc.ctx, tc.username, tc.bid.Id, tc.status)

			tc.expectations(t, r, err)
		})
	}
}
