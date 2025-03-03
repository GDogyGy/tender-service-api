package fetch_test

import (
	ping "TenderServiceApi/internal/handlers/ping/fetch"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePing(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		prepare      func()
		expectations func(t *testing.T, statusCode int)
	}{
		{
			name:    "Success",
			url:     "/api/ping/",
			method:  http.MethodGet,
			prepare: func() {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusOK, statusCode)
			},
		},
		{
			name:    "Error",
			url:     "/api/bids/my",
			method:  http.MethodPost,
			prepare: func() {},
			expectations: func(t *testing.T, statusCode int) {
				assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.prepare()

			handler := ping.NewHandler()

			r := httptest.NewRequest(tc.method, tc.url, nil)

			w := httptest.NewRecorder()

			handler.Ping(w, r)
			tc.expectations(t, w.Result().StatusCode)
		})
	}
}
