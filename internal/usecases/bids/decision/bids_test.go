package decision_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"

	"TenderServiceApi/internal/model"
	"TenderServiceApi/internal/usecases/bids/decision"
	"github.com/stretchr/testify/assert"
)

func TestUseCaseSubmitDecision(t *testing.T) {
	cases := []struct {
		name           string
		ctx            context.Context
		username       string
		bidID          string
		decision       string
		organizationID string
		prepare        func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify)
		expectations   func(t *testing.T, result model.Bids, err error)
	}{
		{
			name:           "Success bid APPROVED",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "APPROVED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{
					{
						Id:          "3",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "3",
					},
					{
						Id:          "2",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "2",
					},
				}, nil)

				decisionRepository.On("Save", mock.Anything, "1", "APPROVED", "1").Return(nil)

				bidRepository.On("UpdateStatus", mock.Anything, "1", "APPROVED").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "APPROVED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderRepository.On("UpdateStatus", mock.Anything, "1", "CLOSED").Return(model.Tender{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					ServiceType: "Test",
					Status:      "CLOSED",
					Version:     1,
					Responsible: "1",
				}, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, "APPROVED", result.Status)
			},
		},
		{
			name:           "Success bid CANCELED",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{
					{
						Id:          "3",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "3",
					},
					{
						Id:          "2",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "2",
					},
				}, nil)

				decisionRepository.On("Save", mock.Anything, "1", "REJECTED", "1").Return(nil)

				bidRepository.On("UpdateStatus", mock.Anything, "1", "CANCELED").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "CANCELED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.NoError(t, err)
				assert.Equal(t, "CANCELED", result.Status)
			},
		},
		{
			name:           "Error bid organizationVerify CheckResponsible",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{}, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error bidRepository organizationVerify FetchById",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error bidRepository tenderVerify CheckResponsible",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "").Return(false, model.NotFindResponsible)
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error decisionRepository FetchByBidID",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error decisionRepository Save",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{
					{
						Id:          "3",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "3",
					},
					{
						Id:          "2",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "2",
					},
				}, nil)

				decisionRepository.On("Save", mock.Anything, "1", "REJECTED", "1").Return(fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error bid UpdateStatus CANCELED",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "REJECTED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{
					{
						Id:          "3",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "3",
					},
					{
						Id:          "2",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "2",
					},
				}, nil)

				decisionRepository.On("Save", mock.Anything, "1", "REJECTED", "1").Return(nil)

				bidRepository.On("UpdateStatus", mock.Anything, "1", "CANCELED").Return(model.Bids{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error bid UpdateStatus APPROVED",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "APPROVED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{
					{
						Id:          "3",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "3",
					},
					{
						Id:          "2",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "2",
					},
				}, nil)

				decisionRepository.On("Save", mock.Anything, "1", "APPROVED", "1").Return(nil)

				bidRepository.On("UpdateStatus", mock.Anything, "1", "APPROVED").Return(model.Bids{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.Empty(t, result)
				assert.Error(t, err)
			},
		},
		{
			name:           "Error tender UpdateStatus CLOSED",
			ctx:            context.Background(),
			username:       "user1",
			bidID:          "1",
			decision:       "APPROVED",
			organizationID: "1",
			prepare: func(bidRepository *MockbidRepository, decisionRepository *MockdecisionRepository, tenderRepository *MocktenderRepository, organizationVerify *MockuseCaseOrganizationVerify, tenderVerify *MockuseCaseTenderVerify) {
				organizationVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(model.OrganizationResponsible{
					Id:             "1",
					OrganizationId: "1",
					UserId:         "1",
				}, nil)

				bidRepository.On("FetchById", mock.Anything, "1").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderVerify.On("CheckResponsible", mock.Anything, "user1", "1").Return(true, nil)

				decisionRepository.On("FetchByBidID", mock.Anything, "1").Return([]model.BidsDecisions{
					{
						Id:          "3",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "3",
					},
					{
						Id:          "2",
						BidID:       "1",
						Decision:    "APPROVED",
						Responsible: "2",
					},
				}, nil)

				decisionRepository.On("Save", mock.Anything, "1", "APPROVED", "1").Return(nil)

				bidRepository.On("UpdateStatus", mock.Anything, "1", "APPROVED").Return(model.Bids{
					Id:          "1",
					Name:        "Test",
					Description: "Test",
					Status:      "APPROVED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1",
				}, nil)

				tenderRepository.On("UpdateStatus", mock.Anything, "1", "CLOSED").Return(model.Tender{}, fmt.Errorf(""))
			},
			expectations: func(t *testing.T, result model.Bids, err error) {
				assert.NotEmpty(t, result)
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bidRepository := NewMockbidRepository(t)
			decisionRepository := NewMockdecisionRepository(t)
			tenderRepository := NewMocktenderRepository(t)

			useCaseOrganizationVerify := NewMockuseCaseOrganizationVerify(t)
			useCaseTenderVerify := NewMockuseCaseTenderVerify(t)

			tc.prepare(bidRepository, decisionRepository, tenderRepository, useCaseOrganizationVerify, useCaseTenderVerify)

			useCase := decision.NewService(bidRepository, decisionRepository, tenderRepository, useCaseOrganizationVerify, useCaseTenderVerify)
			r, err := useCase.SubmitDecision(tc.ctx, tc.username, tc.bidID, tc.decision, tc.organizationID)

			tc.expectations(t, r, err)
		})
	}
}
