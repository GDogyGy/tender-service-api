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

func TestBidsFeedbackHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	body := []byte(`{
		"description": "Плохое предложение"
	}`)

	bidsID := "0e02782f-7891-4813-86f2-074bb6702ab3"
	user := "user2"

	resp, err := request.Execute(http.MethodPut, "http://localhost:9000/api/bids/"+bidsID+"/feedback?username="+user, body)
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
	result := Db.QueryRowxContext(context.Background(), `SELECT count(*) FROM bid_feedback where description =  'Плохое предложение' and bid_id = $1`, bidsID)
	err = result.Err()
	assert.NoError(t, err)

	var c int
	_ = result.Scan(&c)

	assert.Equal(t, 1, c)
}
