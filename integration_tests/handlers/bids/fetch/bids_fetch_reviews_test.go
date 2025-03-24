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

func TestFetchReviewsHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	user := "user3"
	author := "user2"
	tenderID := "81c19f2e-1ed3-46f9-bce3-aa0206caf30e"
	organizationID := "109dae74-bbf6-484b-8b2b-ab7c97967905"

	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/bids/"+tenderID+"/reviews?organizationId="+organizationID+"&username="+user+"&authorUsername="+author, nil)
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

	var feedbacks []transport.Feedback
	err = json.Unmarshal(body, &feedbacks)
	assert.NoError(t, err)

	assert.Greater(t, len(feedbacks), 0)

	r := Db.QueryRowxContext(context.Background(), `SELECT count(*) FROM bids LEFT JOIN bid_feedback ON bid_feedback.bid_id = bids.id WHERE bids.tender_id = $1 and bids.version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id) and bids.responsible = ANY (SELECT organization_responsible.id FROM organization_responsible LEFT JOIN employee on employee.id = organization_responsible.user_id WHERE employee.username = $2 and organization_responsible.organization_id = $3) and bid_feedback.id is not null`, tenderID, author, organizationID)
	assert.NoError(t, r.Err())

	var count int
	err = r.Scan(&count)
	assert.NotEqual(t, count, 0)
	assert.Equal(t, count, len(feedbacks))

}
