package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ApiResponse[T any] struct {
	StatusCode int
	Data       T
}

type NewApi[T any] struct {
	ctx     *AppContext
	method  string
	url     string
	body    io.Reader
	Request *http.Request
}

// Using the std `http.NewRequest` pattern, `New` instantiates a request
func (m *NewApi[T]) New() (*NewApi[T], error) {
	Info(fmt.Sprintf("Request URL: %s", m.url))
	Info(fmt.Sprintf("Request Method: %s", m.method))
	request, err := http.NewRequest(m.method, m.url, m.body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accepts", "application/json")

	m.Request = request

	return m, err
}

// Run the request instantiated by `New`
func (m *NewApi[T]) makeRequest() (ApiResponse[T], error) {
	if m.ctx.ApiKey != "" {
		bearer := fmt.Sprintf("Bearer %s", m.ctx.ApiKey)
		m.Request.Header.Set("Authorization", bearer)
	}

	var data T
	result := ApiResponse[T]{}
	client := new(http.Client)
	response, err := client.Do(m.Request)
	if err != nil {
		return result, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	result.StatusCode = response.StatusCode
	if response.StatusCode == http.StatusOK {
		json.Unmarshal(body, &data)
		result.Data = data

		return result, err
	}

	// Check the Content-Type header
	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		return result, fmt.Errorf("received HTML error page (status code %d)", response.StatusCode)
	}

	var errorResponse ApiErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		errorBody := string(body)
		return result, fmt.Errorf(errorBody)
	}
	Error(fmt.Sprintf("Response Body: %s", errorResponse.GetError()))
	Error(fmt.Sprintf("Status Code: %d", response.StatusCode))
	return result, errorResponse.GetError()
}

func (m *NewApi[T]) Do() (ApiResponse[T], error) {
	response, err := m.makeRequest()
	if err != nil && response.StatusCode == http.StatusUnauthorized {
		Error("Session expired")
		err := RefreshToken(m.ctx)
		if err != nil {
			return ApiResponse[T]{}, err
		}
		return m.makeRequest()
	}

	return response, err
}

// Convenience methods
func NewApiGet[T any](ctx *AppContext, url string) (*NewApi[T], error) {
	config := NewApi[T]{ctx: ctx, method: "GET", url: url}
	request, err := config.New()
	return request, err
}

func NewApiPost[T any](ctx *AppContext, url string, body io.Reader) (*NewApi[T], error) {
	config := NewApi[T]{ctx: ctx, method: "POST", url: url, body: body}
	request, err := config.New()
	return request, err
}

func NewApiPut[T any](ctx *AppContext, url string, body io.Reader) (*NewApi[T], error) {
	config := NewApi[T]{ctx: ctx, method: "PUT", url: url, body: body}
	request, err := config.New()
	return request, err
}

func NewApiDelete[T any](ctx *AppContext, url string) (*NewApi[T], error) {
	config := NewApi[T]{ctx: ctx, method: "DELETE", url: url}
	request, err := config.New()
	return request, err
}

// API error parsing
type ApiErrorResponse struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_description"`
	Message      string `json:"message"`
	Results      string `json:"results"`
	Details      struct {
		Item struct {
			Message string `json:"message"`
		} `json:"d0"`
	} `json:"details"`
}

func (e *ApiErrorResponse) GetError() error {
	var result []string
	if e.Error != "" {
		result = append(result, e.Error)
	}
	if e.Message != "" {
		result = append(result, e.Message)
	}
	if e.Results != "" {
		result = append(result, e.Results)
	}
	if e.Details.Item.Message != "" {
		result = append(result, e.Details.Item.Message)
	}
	if e.ErrorMessage != "" {
		result = append(result, e.ErrorMessage)
	}
	if len(result) == 0 {
		result = append(result, "something broke")
	}

	fullResult := strings.Join(result, "\n")
	return fmt.Errorf(fullResult)
}
