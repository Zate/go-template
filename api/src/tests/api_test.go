package tests

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// struct for building a NewRequest
type request struct {
	method   string
	endpoint string
	body     io.Reader
	headers  map[string]string
}

// struct for a test which has the request and the expected response
type test struct {
	name        string
	request     request
	expected    string
	respCode    int
	errorString string
	statusOnly  bool
}

func TestAPI(t *testing.T) {
	client := &http.Client{}
	// create a slice of tests
	tests := []test{
		{
			name: "root",
			request: request{
				method:   "GET",
				endpoint: "/",
				body:     nil,
			},
			respCode:   http.StatusOK,
			statusOnly: true,
		},
		{
			name: "Healthcheck",
			request: request{
				method:   "GET",
				endpoint: "/healthcheck",
				body:     nil,
			},
			respCode:   http.StatusOK,
			statusOnly: true,
		},
		{
			name: "JSON Status",
			request: request{
				method:   "GET",
				endpoint: "/status",
				body:     nil,
			},
			statusOnly: true,
			respCode:   http.StatusOK,
		},
		{
			name: "Debug",
			request: request{
				method:   "GET",
				endpoint: "/debug",
				body:     nil,
			},
			respCode:   http.StatusOK,
			statusOnly: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make the API request and get the response...
			req, err := http.NewRequest(tt.request.method, fmt.Sprintf("http://%s:8080%s", os.Getenv("DOCKER_GATEWAY_HOST"), tt.request.endpoint), tt.request.body)
			if err != nil {
				t.Fatal(err)
			}

			if tt.request.headers != nil {
				for k, v := range tt.request.headers {
					req.Header.Set(k, v)
				}
			}

			resp, err := client.Do(req)

			if err != nil {
				t.Fatal(err)
			}

			// Check the response status code
			assert.Equal(t, tt.respCode, resp.StatusCode, "they should be equal")

			if tt.statusOnly {
				return
			}

			// Read the response body
			body, _ := io.ReadAll(resp.Body)

			// Check the response body
			assert.Equal(t, tt.expected, strings.TrimSpace(string(body)), "they should be equal")
		})
	}
}
