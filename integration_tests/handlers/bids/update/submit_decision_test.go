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

func TestSubmitDecisionHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}
	user := "user3"
	bidID := "385a485c-c977-429d-88ee-72fec6aa0444"
	decision := "APPROVED"
	organizationId := "109dae74-bbf6-484b-8b2b-ab7c97967905"

	res, err := request.Execute(http.MethodPut, "http://localhost:9000/api/bids/"+bidID+"/submit_decision?username="+user+"&organizationId="+organizationId+"&decision="+decision, []byte{})
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

	var bid transport.Bid
	err = json.Unmarshal(body, &bid)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT bids.status FROM bids  WHERE bids.id = $1 ORDER BY version DESC LIMIT 1`, bidID)
	assert.NoError(t, r.Err())

	var status string
	err = r.Scan(&status)

	assert.NoError(t, err)
	assert.NotEqual(t, "PUBLISHED", bid.Status)
	assert.Equal(t, decision, string(bid.Status))

	r = Db.QueryRowxContext(context.Background(), `SELECT tender.status FROM tender  WHERE tender.id = $1 ORDER BY version desc limit 1`, bid.TenderId)
	assert.NoError(t, r.Err())
	var tenderStatus string
	err = r.Scan(&tenderStatus)
	assert.Equal(t, "CLOSED", tenderStatus)

}
