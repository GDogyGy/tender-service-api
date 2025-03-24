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

func TestFetchListByTenderHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	tenderID := "81c19f2e-1ed3-46f9-bce3-aa0206caf30e"
	user := "user2"
	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/bids/"+tenderID+"/list?username="+user, nil)
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

	var bids []transport.Bid
	err = json.Unmarshal(body, &bids)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT count(*) FROM bids WHERE tender_id = $1 and id = $2`, tenderID, bids[0].Id)
	assert.NoError(t, r.Err())

	var count int
	err = r.Scan(&count)
	assert.NoError(t, err)
	assert.Greater(t, len(bids), 0)
	assert.NotEqual(t, count, 0)
}
