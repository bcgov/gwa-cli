package pkg

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestConvertVersion(t *testing.T) {
	assert.Equal(t, 10200, convertVersion("v1.2.0"))
	assert.Equal(t, 101905, convertVersion("v10.19.5"))
}

func TestCompareVersion(t *testing.T) {
	assert.True(t, isOutdated("v1.2.0", "v10.3.0"))
	assert.False(t, isOutdated("v1.3.0", "v1.3.0"))
}

func TestCheckForVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", releaseUrl, func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"tag_name": "v1.2.0",
		})
	})
	out := capturer.CaptureStdout(func() {
		ctx := AppContext{
			Version: "v1.1.1",
		}
		CheckForVersion(&ctx)
	})
	assert.Contains(t, out, "A new version (v1.2.0) is available.")
	empty := capturer.CaptureStdout(func() {
		ctx := AppContext{
			Version: "v1.2.0",
		}
		CheckForVersion(&ctx)
	})
	assert.Equal(t, empty, "")
}
