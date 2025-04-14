//go:build integration

package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"TenderServiceApi/integration_tests/helpers/db/client"
	"TenderServiceApi/integration_tests/helpers/server/request"
	"TenderServiceApi/internal/handlers/rest/types/transport"
	"github.com/stretchr/testify/assert"
)

func TestFetchListByUserHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/tenders/my?username=user1", nil)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	defer func() { _ = res.Body.Close() }()

	assert.NotEmpty(t, body)

	var tenders []transport.Tender
	err = json.Unmarshal(body, &tenders)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT t.name FROM tender t LEFT JOIN organization_responsible o ON t.responsible = o.organization_id LEFT JOIN employee e ON e.id = o.user_id WHERE e.username = $1  LIMIT 1`, "user1")
	assert.NoError(t, r.Err())

	var searchName string
	err = r.Scan(&searchName)
	assert.NoError(t, err)

	assert.Equal(t, searchName, tenders[0].Name)
	assert.Greater(t, len(tenders), 0)
}
