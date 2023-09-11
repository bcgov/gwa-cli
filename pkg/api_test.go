package pkg

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type BasicResponse struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
	Id    int    `json:"id"`
}

var ctx = &AppContext{
	ApiKey: "123123123",
}

const URL = "https://test.app"

func TestApiGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", URL, func(r *http.Request) (*http.Response, error) {
		assert.Nil(t, r.Body)
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"name":  "Hello",
			"total": 42,
			"id":    1,
		})
	})

	ctx := &AppContext{
		ApiKey: "123123123",
	}

	r, err := NewApiGet[BasicResponse](ctx, URL)
	if err != nil {
		t.Fatal(err)
	}
	response, err := r.Do()

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, BasicResponse{Name: "Hello", Total: 42, Id: 1}, response.Data)
	assert.Nil(t, err)
}

func TestApiPost(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", URL, func(r *http.Request) (*http.Response, error) {
		assert.NotNil(t, r.PostForm)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"), "content-type is set to form")
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"name":  "Hello",
			"total": 42,
			"id":    1,
		})
	})

	data := make(url.Values)
	data.Set("name", "Hello")
	data.Set("total", "42")

	r, err := NewApiPost[BasicResponse](ctx, URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Request.PostForm = data
	r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := r.Do()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, BasicResponse{Name: "Hello", Total: 42, Id: 1}, response.Data)
	assert.Nil(t, err)
}

func TestApiPut(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", URL, func(r *http.Request) (*http.Response, error) {
		assert.NotNil(t, r.PostForm)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"), "content-type is set to form")
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"name":  "World",
			"total": 42,
			"id":    1,
		})
	})

	data := make(url.Values)
	data.Set("name", "World")
	data.Set("total", "42")

	r, err := NewApiPut[BasicResponse](ctx, URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Request.PostForm = data
	r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := r.Do()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, BasicResponse{Name: "World", Total: 42, Id: 1}, response.Data)
	// assert.Nil(t, err)
}

func TestApiDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", URL, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"name":  "World",
			"total": 42,
			"id":    1,
		})
	})

	r, err := NewApiDelete[BasicResponse](ctx, URL)
	if err != nil {
		t.Fatal(err)
	}
	response, err := r.Do()

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, BasicResponse{Name: "World", Total: 42, Id: 1}, response.Data)
	assert.Nil(t, err)
}

func TestApiError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", URL, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(500, map[string]interface{}{
			"error":             "Service unavailable",
			"error_description": "The server is not responding",
		})
	})

	r, err := NewApiGet[BasicResponse](ctx, URL)
	if err != nil {
		t.Fatal(err)
	}
	response, err := r.Do()

	assert.NotNil(t, err)
	assert.Equal(t, 500, response.StatusCode)
}

func TestEmptyApiError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", URL, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, ""), nil
	})

	r, err := NewApiGet[BasicResponse](ctx, URL)
	if err != nil {
		t.Fatal(err)
	}
	response, err := r.Do()

	assert.NotNil(t, err)
	assert.Equal(t, 500, response.StatusCode)
}

func TestErrorStruct(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		payload  string
	}{
		{
			name:     "standard error",
			expected: "Service Unavailable\nThe server is not responding",
			payload:  `{"error": "Service Unavailable", "error_description":  "The server is not responding"}`,
		},
		{
			name:     "dense error",
			expected: "jwt invalid\nyou are not authorized to visit",
			payload:  `{"message": "jwt invalid", "details": {"d0": {"message": "you are not authorized to visit"}}}`,
		},
		{
			name:     "error only",
			expected: "Service Unavailable",
			payload:  `{"error": "Service Unavailable"}`,
		},
		{
			name:     "description only",
			expected: "The server is not responding",
			payload:  `{"error_description":  "The server is not responding"}`,
		},
		{
			name:     "empty payload",
			expected: "something broke",
			payload:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result ApiErrorResponse
			json.Unmarshal([]byte(tt.payload), &result)

			assert.ErrorContains(t, result.GetError(), tt.expected)
		})
	}
}
