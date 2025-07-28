//go:build integration

package test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"TenderServiceApi/integration_tests/helpers/db/client"
	"TenderServiceApi/integration_tests/helpers/server/request"
	"github.com/stretchr/testify/assert"
)

func TestTenderCreateHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	body := []byte(`{
			"name": "Предложение 5",
			"description": "Проверю КАЧЕСТВЕННО !!! рестораны на качество услуг",
			"status": "CREATED",
			"tenderId": "81c19f2e-1ed3-46f9-bce3-aa0206caf30e",
			"version": 1,
			"organizationId": "550e8400-e29b-41d4-a716-446655440000",
			"creatorUsername": "user1"
		}`)

	resp, err := request.Execute(http.MethodPost, "http://localhost:9000/api/bids/new", body)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	defer func() { _ = resp.Body.Close() }()

	assert.NotEmpty(t, body)
	result := Db.QueryRowxContext(context.Background(), `SELECT count(*) FROM bids where name =  'Предложение 5'`)
	err = result.Err()
	assert.NoError(t, err)

	var c int
	_ = result.Scan(&c)

	assert.Equal(t, 1, c)
}
