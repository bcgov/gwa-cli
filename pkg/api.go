package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ApiErrorResponse struct {
	Message string `json:"message"`
	Details struct {
		Item struct {
			Message string `json:"message"`
		} `json:"d0"`
	} `json:"details"`
}

type ApiResponse[T any] struct {
	StatusCode int
	Data       T
}

func Api[T any](ctx *AppContext, url string, method string, requestBody io.Reader) (ApiResponse[T], error) {
	var data T
	result := ApiResponse[T]{}
	client := http.Client{}
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return result, err
	}
	bearer := fmt.Sprintf("Bearer %s", ctx.ApiKey)
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
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
		json.Unmarshal(body, &errorResponse)
		errMessage := strings.Join([]string{errorResponse.Message, errorResponse.Details.Item.Message}, " ")
		return result, errors.New(errMessage)
	}
}

func ApiGet[T any](ctx *AppContext, url string) (ApiResponse[T], error) {
	return Api[T](ctx, url, http.MethodGet, nil)
}

func ApiPost[T any](ctx *AppContext, url string, body io.Reader) (ApiResponse[T], error) {
	return Api[T](ctx, url, http.MethodPost, body)
}

func ApiPut[T any](ctx *AppContext, url string, body io.Reader) (ApiResponse[T], error) {
	return Api[T](ctx, url, http.MethodPut, body)
}

// func ApiDelete(ctx *AppContext, url string) (bool, error) {
// 	return Api[bool](ctx, url, http.MethodDelete, nil)
// }
