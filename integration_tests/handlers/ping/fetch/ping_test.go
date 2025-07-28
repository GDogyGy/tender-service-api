//go:build integration

package test

import (
	"net/http"
	"testing"

	"TenderServiceApi/integration_tests/helpers/server/request"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	res, err := request.Execute(http.MethodGet, "http://localhost:9000/api/ping", nil)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
