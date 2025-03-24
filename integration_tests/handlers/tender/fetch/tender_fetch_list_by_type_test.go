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
	"TenderServiceApi/internal/handlers/types/transport"
	"github.com/stretchr/testify/assert"
)

func TestFetchListByTypeHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/tenders?servicetype=Development", nil)
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

	r := Db.QueryRowxContext(context.Background(), `SELECT name FROM tender WHERE service_type = $1 LIMIT 1`, "Development")
	assert.NoError(t, r.Err())

	var searchName string
	err = r.Scan(&searchName)
	assert.NoError(t, err)

	assert.Equal(t, searchName, tenders[0].Name)
	assert.Greater(t, len(tenders), 0)
}
