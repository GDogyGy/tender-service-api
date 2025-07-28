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
	"github.com/stretchr/testify/assert"
)

func TestFetchStatusHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	user := "user2"
	bidID := "385a485c-c977-429d-88ee-72fec6aa0444"

	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/bids/"+bidID+"/status?username="+user, nil)
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

	var status string
	err = json.Unmarshal(body, &status)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT b.status FROM bids b LEFT JOIN organization_responsible o ON b.responsible = o.id LEFT JOIN employee e ON e.id = o.user_id WHERE e.username = $1 AND b.id = $2  LIMIT 1`, user, bidID)
	assert.NoError(t, r.Err())

	var searchStatus string
	err = r.Scan(&searchStatus)
	assert.NoError(t, err)

	assert.Equal(t, searchStatus, status)
}
