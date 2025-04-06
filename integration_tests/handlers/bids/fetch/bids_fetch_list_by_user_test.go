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

func TestFetchListByUserHandler(t *testing.T) {
	if err := client.DeleteAllCollections(context.Background(), Db); err != nil {
		t.Error(err)
	}

	if err := client.PrepareDB(context.Background(), Db); err != nil {
		t.Error(err)
	}

	user := "user1"
	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/bids/my?username="+user, nil)
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

	r := Db.QueryRowxContext(context.Background(), `SELECT count(*) FROM bids left join organization_responsible o on responsible = o.id left join employee e on o.user_id = e.id WHERE e.username = $1 and version = (SELECT MAX(version) FROM bids t2 WHERE bids.id = t2.id)`, user)
	assert.NoError(t, r.Err())

	var count int
	err = r.Scan(&count)
	assert.NoError(t, err)
	assert.Greater(t, len(bids), 0)
	assert.NotEqual(t, count, 0)
	assert.Equal(t, count, len(bids))
}
