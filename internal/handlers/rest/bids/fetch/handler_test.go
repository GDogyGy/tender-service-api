package fetch_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	bids "TenderServiceApi/internal/handlers/bids/fetch"
	"TenderServiceApi/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleFetchListByUser(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/bids/my?username=user1",
			method: http.MethodGet,
			prepare: func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog) {
				useCaseBids.On("FetchListByUser", mock.Anything, "user1").Return([]model.Bids{{
					Id:          "123",
					Name:        "Name",
					Description: "Description",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1b3dd29a-ba01-4374-a79f-0c7b654bea67",
				},
				}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusOK, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/my",
			method:  http.MethodPost,
			prepare: func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/my?username=",
			method:  http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusBadRequest, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/my?username=user1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {
				service.On("FetchListByUser", mock.Anything, "user1").Return([]model.Bids{}, fmt.Errorf(""))
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			},
		},
		{
			name:   "No content",
			url:    "/api/bids/my?username=user1",
			method: http.MethodGet,
			prepare: func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog) {
				useCaseBids.On("FetchListByUser", mock.Anything, "user1").Return([]model.Bids{}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusNotFound, statusCode)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCaseBidsFetch := NewMockuseCaseBidsFetch(t)
			useCaseBidFeedbackFetch := NewMockuseCaseBidFeedbackFetch(t)
			logMock := NewMocklog(t)

			tc.prepare(useCaseBidsFetch, logMock)

			handler := bids.NewHandler(logMock, useCaseBidsFetch, useCaseBidFeedbackFetch)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchListByUser(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}

func TestHandleFetchListByTender(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/bids/1/list?username=user1",
			method: http.MethodGet,
			prepare: func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog) {
				useCaseBids.On("FetchListByTender", mock.Anything, "user1", "1").Return([]model.Bids{{
					Id:          "1",
					Name:        "Name",
					Description: "Description",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1b3dd29a-ba01-4374-a79f-0c7b654bea67",
				},
				}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusOK, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/1/list?username=user1",
			method:  http.MethodPost,
			prepare: func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/1/list?username=",
			method:  http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusBadRequest, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/list?username=user1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {
				service.On("FetchListByTender", mock.Anything, "user1", "1").Return([]model.Bids{}, model.NotFindResponsible)
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusForbidden, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/list?username=user1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {
				service.On("FetchListByTender", mock.Anything, "user1", "1").Return([]model.Bids{}, model.NotFound)
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusNoContent, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/list?username=user1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {
				service.On("FetchListByTender", mock.Anything, "user1", "1").Return([]model.Bids{}, fmt.Errorf(""))
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			},
		},
		{
			name:   "No content",
			url:    "/api/bids/1/list?username=user1",
			method: http.MethodGet,
			prepare: func(useCaseBids *MockuseCaseBidsFetch, log *Mocklog) {
				useCaseBids.On("FetchListByTender", mock.Anything, "user1", "1").Return([]model.Bids{}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusNoContent, statusCode)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCaseBidsFetch := NewMockuseCaseBidsFetch(t)
			useCaseBidFeedbackFetch := NewMockuseCaseBidFeedbackFetch(t)
			logMock := NewMocklog(t)

			tc.prepare(useCaseBidsFetch, logMock)

			handler := bids.NewHandler(logMock, useCaseBidsFetch, useCaseBidFeedbackFetch)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchListByTender(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}

func TestHandleFetchStatus(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(service *MockuseCaseBidsFetch, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/bids/1/status?username=username",
			method: http.MethodGet,
			prepare: func(useCase *MockuseCaseBidsFetch, log *Mocklog) {
				useCase.On("FetchStatus", mock.Anything, "username", "1").Return(model.Bids{
					Id:          "1",
					Name:        "Name",
					Description: "Description",
					Status:      "PUBLISHED",
					TenderId:    "1",
					Version:     1,
					Responsible: "1b3dd29a-ba01-4374-a79f-0c7b654bea67",
				}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusOK, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/status?username=username",
			method: http.MethodGet,
			prepare: func(useCase *MockuseCaseBidsFetch, log *Mocklog) {
				useCase.On("FetchStatus", mock.Anything, "username", "1").Return(model.Bids{}, model.NotFindResponsible)
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusForbidden, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/status?username=username",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {
				service.On("FetchStatus", mock.Anything, "username", "1").Return(model.Bids{}, sql.ErrNoRows)
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusForbidden, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/status?username=username",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {
				service.On("FetchStatus", mock.Anything, "username", "1").Return(model.Bids{}, fmt.Errorf(""))
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/1/status?username=username",
			method:  http.MethodPost,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/1/status?username=",
			method:  http.MethodGet,
			prepare: func(service *MockuseCaseBidsFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusBadRequest, statusCode)
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCaseBidsFetch := NewMockuseCaseBidsFetch(t)
			useCaseBidFeedbackFetch := NewMockuseCaseBidFeedbackFetch(t)
			logMock := NewMocklog(t)

			tc.prepare(useCaseBidsFetch, logMock)

			handler := bids.NewHandler(logMock, useCaseBidsFetch, useCaseBidFeedbackFetch)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchStatus(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}

func TestHandleFetchReviews(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(service *MockuseCaseBidFeedbackFetch, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/bids/2/reviews?username=username1&authorUsername=username2&organizationId=1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidFeedbackFetch, log *Mocklog) {
				service.On("FetchReviews", mock.Anything, "username1", "2", "username2", "1").Return([]model.BidFeedback{
					{
						Id:          "1",
						BidID:       "1",
						Description: "Description",
						Responsible: "1b3dd29a-ba01-4374-a79f-0c7b654bea67",
						CreatedAt:   "2025-02-21 10:17:30.541167",
					},
				}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusOK, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/reviews?username=username&authorUsername=username2&organizationId=1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidFeedbackFetch, log *Mocklog) {
				service.On("FetchReviews", mock.Anything, "username", "1", "username2", "1").Return([]model.BidFeedback{}, model.NotFindResponsible)
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusForbidden, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/reviews?username=username&authorUsername=username2&organizationId=1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidFeedbackFetch, log *Mocklog) {
				service.On("FetchReviews", mock.Anything, "username", "1", "username2", "1").Return([]model.BidFeedback{}, sql.ErrNoRows)
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusForbidden, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/bids/1/reviews?username=username&authorUsername=username2&organizationId=1",
			method: http.MethodGet,
			prepare: func(service *MockuseCaseBidFeedbackFetch, log *Mocklog) {
				service.On("FetchReviews", mock.Anything, "username", "1", "username2", "1").Return([]model.BidFeedback{}, fmt.Errorf(""))
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/1/reviews?username=username&authorUsername=username2&organizationId=1",
			method:  http.MethodPost,
			prepare: func(service *MockuseCaseBidFeedbackFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/1/reviews?username&authorUsername&organizationId=1",
			method:  http.MethodGet,
			prepare: func(service *MockuseCaseBidFeedbackFetch, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusBadRequest, statusCode)
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCaseBidsFetch := NewMockuseCaseBidsFetch(t)
			useCaseBidFeedbackFetch := NewMockuseCaseBidFeedbackFetch(t)
			logMock := NewMocklog(t)

			tc.prepare(useCaseBidFeedbackFetch, logMock)

			handler := bids.NewHandler(logMock, useCaseBidsFetch, useCaseBidFeedbackFetch)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchReviews(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}
