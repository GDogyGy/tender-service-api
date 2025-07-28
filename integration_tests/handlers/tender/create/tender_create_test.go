//go:build integration

package test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"TenderServiceApi/integration_tests/helpers/db/client"
	"TenderServiceApi/integration_tests/helpers/server/request"
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

	res, err := request.Execute(http.MethodPost, "http://localhost:9000/api/tenders/new", body)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)

	assert.Equal(t, 200, res.StatusCode)
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	defer func() { _ = res.Body.Close() }()

	expectationYo(t, Db, body)
}

func expectationYo(t *testing.T, Db *sqlx.DB, body []byte) {
	assert.NotEmpty(t, body)
	res := Db.QueryRowxContext(context.Background(), `SELECT count(*) FROM tender where name =  'Tender: new obj'`)
	err := res.Err()
	assert.NoError(t, err)
	var c int
	_ = res.Scan(&c)

	assert.Equal(t, 1, c)

}
