package request

import (
	"bytes"
	"net/http"
	"time"
)

func Execute(method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	clientHttp := &http.Client{
		Timeout: 30 * time.Second,
	}

	return clientHttp.Do(req)
}
