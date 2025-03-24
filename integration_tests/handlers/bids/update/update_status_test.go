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

func TestUpdateStatusHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}
	user := "user1"
	bidID := "f28798c2-314b-4ba6-bb49-9ec23d8a2f65"
	statusNew := "PUBLISHED"

	res, err := request.Execute(http.MethodPut, "http://localhost:9000/api/bids/"+bidID+"/status?username="+user+"&status="+statusNew, []byte{})
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

	var tender transport.Tender
	err = json.Unmarshal(body, &tender)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT bids.status FROM bids  WHERE bids.id = $1 ORDER BY version DESC LIMIT 1`, bidID)
	assert.NoError(t, r.Err())

	var status string
	err = r.Scan(&status)

	assert.NoError(t, err)
	assert.NotEmpty(t, status)
	assert.NotEqual(t, "CREATED", status)
	assert.Equal(t, statusNew, status)
}
