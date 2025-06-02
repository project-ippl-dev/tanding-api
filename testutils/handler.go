package testutils

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHttpRequestParam struct {
	HttpMethod string
	Url        string
	ReqBody    interface{}
	ReqRawBody *bytes.Buffer
	ReqHeaders map[string]string
}

func MockHttpRequest(t *testing.T, input MockHttpRequestParam) (rr *httptest.ResponseRecorder, httpReq *http.Request) {
	if input.ReqBody != nil {
		req, err := json.Marshal(input.ReqBody)
		if err != nil {
			t.Fatalf("fail marshal reqBody: %s", err.Error())
		}
		httpBody := bytes.NewBuffer(req)
		httpReq = httptest.NewRequest(input.HttpMethod, input.Url, httpBody)
	} else if input.ReqRawBody != nil {
		httpReq = httptest.NewRequest(input.HttpMethod, input.Url, input.ReqRawBody)
	} else {
		httpReq = httptest.NewRequest(input.HttpMethod, input.Url, nil)
	}

	for key, val := range input.ReqHeaders {
		httpReq.Header.Add(key, val)
	}

	if httpReq.Header.Get(echo.HeaderContentType) == "" {
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	rr = httptest.NewRecorder()
	return
}
