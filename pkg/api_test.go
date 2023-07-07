package pkg

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type BasicResponse struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

func TestApiGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://test.app", func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"name":  "Hello",
			"total": 42,
		})
	})

	ctx := &AppContext{
		ApiKey: "123123123",
	}

	response, err := Api[BasicResponse](ctx, "https://test.app", "GET", nil)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, BasicResponse{Name: "Hello", Total: 42}, response.Data)
	assert.Nil(t, err)
}
