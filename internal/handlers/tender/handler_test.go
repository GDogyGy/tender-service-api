package tender_test

import (
	"TenderServiceApi/internal/handlers/tender"
	"TenderServiceApi/internal/model"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleTenderFetchList(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(service *Mockservice, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/tenders?servicetype=11",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("FetchList", mock.Anything, "11").Return([]model.Tender{{
					Id:          "123",
					Name:        "Name",
					Description: "Description",
					ServiceType: "Development",
					Status:      "PUBLISHED",
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
			url:     "/api/tenders",
			method:  http.MethodPost,
			prepare: func(service *Mockservice, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/tenders?servicetype=",
			method:  http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusBadRequest, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/tenders?servicetype=Development",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("FetchList", mock.Anything, "Development").Return([]model.Tender{}, fmt.Errorf(""))
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/tenders",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("FetchList", mock.Anything, "").Return([]model.Tender{}, nil)
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
			serviceMock := NewMockservice(t)
			logMock := NewMocklog(t)
			organizationResponsibleFacade := NewMockorganizationResponsibleFacade(t)

			tc.prepare(serviceMock, logMock)

			handler := tender.NewHandler(logMock, serviceMock, organizationResponsibleFacade)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchList(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}

func TestHandleTenderFetchListByUser(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(service *Mockservice, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/tenders/my?username=username",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("FetchListByUser", mock.Anything, "username").Return([]model.Tender{{
					Id:          "123",
					Name:        "Name",
					Description: "Description",
					ServiceType: "Development",
					Status:      "PUBLISHED",
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
			url:     "/api/tenders/my?username=username",
			method:  http.MethodPost,
			prepare: func(service *Mockservice, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/tenders/my?username",
			method:  http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusBadRequest, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/tenders/my?username=username",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("FetchListByUser", mock.Anything, "username").Return([]model.Tender{}, errors.New(""))
				log.On("Error", mock.Anything).Return("")
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusInternalServerError, statusCode)
			},
		},
		{
			name:   "Error",
			url:    "/api/tenders/my?username=username",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("FetchListByUser", mock.Anything, "username").Return([]model.Tender{}, nil)
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
			serviceMock := NewMockservice(t)
			logMock := NewMocklog(t)
			organizationResponsibleFacade := NewMockorganizationResponsibleFacade(t)

			tc.prepare(serviceMock, logMock)

			handler := tender.NewHandler(logMock, serviceMock, organizationResponsibleFacade)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchListByUser(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}

func TestHandleTenderTenderStatus(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func(service *Mockservice, log *Mocklog)
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:   "Success",
			url:    "/api/tenders/status?username=username&tenderId=1",
			method: http.MethodGet,
			prepare: func(service *Mockservice, log *Mocklog) {
				service.On("CheckResponsibleTender", mock.Anything, "username", "1").Return(true, nil)
				service.On("FetchTenderById", mock.Anything, "1").Return(model.Tender{
					Id:          "1",
					Name:        "Name",
					Description: "Description",
					ServiceType: "Development",
					Status:      "PUBLISHED",
					Responsible: "1b3dd29a-ba01-4374-a79f-0c7b654bea67",
				}, nil)
				service.On("FetchTenderStatus", mock.Anything, "username", "1").Return(model.Tender{
					Id:          "1",
					Name:        "Name",
					Description: "Description",
					ServiceType: "Development",
					Status:      "PUBLISHED",
					Responsible: "1b3dd29a-ba01-4374-a79f-0c7b654bea67",
				}, nil)
			},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusOK, statusCode)
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			serviceMock := NewMockservice(t)
			logMock := NewMocklog(t)
			organizationResponsibleFacade := NewMockorganizationResponsibleFacade(t)

			tc.prepare(serviceMock, logMock)

			handler := tender.NewHandler(logMock, serviceMock, organizationResponsibleFacade)

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.FetchTenderStatus(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}
