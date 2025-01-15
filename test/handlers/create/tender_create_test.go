//go:build integration

package test

import (
	"TenderServiceApi/internal/handlers/tender"
	"bytes"
	"net/http"
	"net/http/httptest"
)

// TODO: получается без makefile не поднять тесты интеграционные. Спросить у димы валидно ли если докер поднимает тестовую базу только через makefile
func (s *APITestSuite) TestTenderCreateHandler() {
	url := "/api/tenders/new"

	body := []byte(`{
		"name": "Tender: inspect qualification",
		"description": "Проверить квалификацию сотрудников",
		"serviceType": "Examination",
		"status": "CLOSED",
		"organizationId": "0577781b-f009-4298-b4cb-ffa17893d6c3",
		"creatorUsername": "user1"
	}`)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		s.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tender.Handler{log * Mocklog}.Create)

	handler.ServeHTTP(rr, req)

	s.Equal(200, rr.Code)
}
