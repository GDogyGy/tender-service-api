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

func TestRollbackHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}
	user := "user3"
	bidID := "0e02782f-7891-4813-86f2-074bb6702ab3"
	body := []byte(`{
			"name": "Bid: edit",
			"description": "Bid edit success"
		}`)
	res, err := request.Execute(http.MethodPatch, "http://localhost:9000/api/bids/"+bidID+"/edit?username="+user, body)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	res, err = request.Execute(http.MethodPut, "http://localhost:9000/api/bids/"+bidID+"/rollback/1?username="+user, body)
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

	assert.NotEmpty(t, body)

	var bid transport.Bid
	err = json.Unmarshal(body, &bid)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT name, max(bids.version) as version FROM bids  WHERE bids.id = $1 GROUP BY name ORDER BY version DESC  LIMIT 1`, bidID)
	assert.NoError(t, r.Err())

	var searchName string
	var version int
	err = r.Scan(&searchName, &version)

	assert.NoError(t, err)

	assert.NotEqual(t, searchName, "Bid: edit")
	assert.Equal(t, version, bid.Version)
}
