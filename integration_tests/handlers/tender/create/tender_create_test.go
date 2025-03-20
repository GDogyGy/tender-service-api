//go:build integration

package test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"TenderServiceApi/integration_tests/db/client"
	"github.com/jmoiron/sqlx"
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
			"name": "Tender: new obj",
			"description": "Проверить квалификацию сотрудников",
			"serviceType": "Examination",
			"status": "CREATED",
			"organizationId": "550e8400-e29b-41d4-a716-446655440000",
			"creatorUsername": "user1"
	}`)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:9000/api/tenders/new", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}

	clientHttp := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := clientHttp.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer func() { _ = res.Body.Close() }()
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body, err = io.ReadAll(res.Body)

	if err != nil {
		t.Error(err)
	}

	expectationYo(t, Db, body)
}

func expectationYo(t *testing.T, Db *sqlx.DB, body []byte) {
	assert.NotEmpty(t, body)
	res := Db.QueryRowxContext(context.Background(), `SELECT FROM tender where name =  'Tender: new obj'`)
	err := res.Err()

	assert.NoError(t, err)
}
