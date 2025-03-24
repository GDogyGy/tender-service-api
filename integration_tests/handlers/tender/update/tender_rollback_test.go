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
	user := "user1"
	tenderID := "cc1f2b36-c45e-43c8-bca7-51112a43171b"
	body := []byte(`{
			"name": "Tender: edit",
			"description": "edit success"
		}`)
	res, err := request.Execute(http.MethodPatch, "http://localhost:9000/api/tenders/"+tenderID+"/edit?username="+user, body)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	res, err = request.Execute(http.MethodPut, "http://localhost:9000/api/tenders/"+tenderID+"/rollback/1?username="+user, body)
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

	var tender transport.Tender
	err = json.Unmarshal(body, &tender)
	assert.NoError(t, err)

	r := Db.QueryRowxContext(context.Background(), `SELECT name, max(tender.version) as version FROM tender  WHERE tender.id = $1 GROUP BY name ORDER BY version DESC  LIMIT 1`, tenderID)
	assert.NoError(t, r.Err())

	var searchName string
	var version int
	err = r.Scan(&searchName, &version)

	assert.NoError(t, err)

	assert.NotEqual(t, searchName, "Tender: edit")
	assert.Equal(t, version, tender.Version)
}
