package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPParams struct {
	URL      string
	Body     io.Reader
	Method   string //HTTP Method
	Headers  []HTTPHeader
	Response interface{}
}

type HTTPHeader struct {
	Key   string
	Value string
}

func HTTPClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second, CheckRedirect: checkRedirect}
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	var err error
	for key, val := range via[0].Header {
		req.Header[key] = val
	}
	return err
}

func HTTPRequest(arg HTTPParams) (statusCode int, err error) {
	httpReq, err := http.NewRequest(arg.Method, arg.URL, arg.Body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in generate new http request : %s", err.Error())
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if arg.Headers != nil {
		for _, header := range arg.Headers {
			httpReq.Header.Set(header.Key, header.Value)
		}
	}
	client := HTTPClient()
	resp, err := client.Do(httpReq)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("something when wrong in execute http request : %s", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in converting response.Body to buffer : %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("error response body : %s", string(body))
	}

	if err := resp.Body.Close(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in close response.Body : %s", err.Error())
	}

	if err := json.Unmarshal(body, &arg.Response); err != nil {
		return http.StatusInternalServerError, fmt.Errorf(" error in unmarshal buffer to struct : %s", err.Error())
	}

	return http.StatusOK, nil
}
