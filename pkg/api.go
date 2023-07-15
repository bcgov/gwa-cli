package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ApiErrorResponse struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_description"`
	Message      string `json:"message"`
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

func (m *NewApi[T]) New() (*NewApi[T], error) {
	request, err := http.NewRequest(m.method, m.url, m.body)
	if err != nil {
		return nil, err
	}
	bearer := fmt.Sprintf("Bearer %s", m.ctx.ApiKey)
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accepts", "application/json")

	m.Request = request

	return m, err
}

func (m *NewApi[T]) Do() (ApiResponse[T], error) {
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
	} else {
		var errorResponse ApiErrorResponse
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return result, fmt.Errorf(string(body))
		}
		return result, errorResponse.GetError()
	}
}

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

// func Api[T any](ctx *AppContext, url string, method string, requestBody io.Reader) (ApiResponse[T], error) {
// 	var data T
// 	result := ApiResponse[T]{}
// 	client := http.Client{}
// 	request, err := http.NewRequest(method, url, requestBody)
// 	if err != nil {
// 		return result, err
// 	}
// 	bearer := fmt.Sprintf("Bearer %s", ctx.ApiKey)
// 	request.Header.Set("Authorization", bearer)
// 	// request.Header.Set("Content-Type", contentType)
// 	request.Header.Set("Accepts", "application/json")
// 	response, err := client.Do(request)
// 	if err != nil {
// 		return result, err
// 	}
//
// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return result, err
// 	}
// 	defer response.Body.Close()
//
// 	result.StatusCode = response.StatusCode
// 	if response.StatusCode == http.StatusOK {
// 		json.Unmarshal(body, &data)
// 		result.Data = data
//
// 		return result, err
// 	} else {
// 		var errorResponse ApiErrorResponse
// 		err := json.Unmarshal(body, &errorResponse)
// 		if err != nil {
// 			return result, fmt.Errorf(string(body))
// 		}
// 		return result, errorResponse.GetError()
// 	}
// }
//
// func ApiGet[T any](ctx *AppContext, url string) (ApiResponse[T], error) {
// 	return Api[T](ctx, url, http.MethodGet, nil)
// }
//
// func ApiPost[T any](ctx *AppContext, url string, body io.Reader) (ApiResponse[T], error) {
// 	return Api[T](ctx, url, http.MethodPost, body)
// }
//
// func ApiPut[T any](ctx *AppContext, url string, body io.Reader) (ApiResponse[T], error) {
// 	return Api[T](ctx, url, http.MethodPut, body)
// }
//
// func ApiDelete[T any](ctx *AppContext, url string) (ApiResponse[T], error) {
// 	return Api[T](ctx, url, http.MethodDelete, nil)
// }
