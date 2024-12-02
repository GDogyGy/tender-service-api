package tender

import (
	"TenderServiceApi/internal/handlers/tender/mocks"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleTenderFetchList(t *testing.T) {
	cases := []struct {
		name      string
		url       string
		param     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			url:   "/api/tenders",
			param: "",
		},
		{
			name:  "Error",
			url:   "/api/tenders",
			param: "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// TODO использую моки для заполнения структуры в NewHandler
			tenderService := mocks.NewService(t)
			tenderLog := mocks.NewLog(t)

			handler := NewHandler(tenderLog, tenderService)

			r := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()

			handler.FetchList(w, r)
			resp := w.Result()
			// TODO: что то не так в моках
			fmt.Println(resp)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
