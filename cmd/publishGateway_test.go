package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

const API_HOST = "myapi.dev"
const configFileContents string = `
_format_version: "1.1"
services:
  - name: Demo_App
    url: /api/demoapp
    plugins: []
`

func TestPublishCommands(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		configFile string
		response   httpmock.Responder
		expect     string
		args       []string
	}{
		{
			name:       "successful straight publish",
			setup:      nil,
			configFile: "config.yaml",
			response:   httpmock.NewStringResponder(200, `{"id": 1}`),
			expect:     "Gateway config published",
			args:       []string{"config.yaml"},
		},
		{
			name:       "api error",
			setup:      nil,
			configFile: "config.yaml",
			response:   httpmock.NewStringResponder(500, `{"error": "something went wrong"}`),
			expect:     "something went wrong",
			args:       []string{"config.yaml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("PUT", "https://"+API_HOST+"/namespaces/ns-sampler/gateway", tt.response)
			cwd := t.TempDir()

			if tt.setup != nil {
				tt.setup()
			}
			if tt.configFile != "" {
				filePath := filepath.Join(cwd, tt.configFile)
				os.WriteFile(filePath, []byte(configFileContents), 0644)
			}
			ctx := &pkg.AppContext{
				Cwd:       cwd,
				ApiHost:   API_HOST,
				Namespace: "ns-sampler",
			}

			args := append([]string{"publish-gateway"}, tt.args...)
			mainCmd := &cobra.Command{
				Use: "gwa",
			}
			mainCmd.AddCommand(NewPublishGatewayCmd(ctx))
			mainCmd.SetArgs(args)
			out := capturer.CaptureOutput(func() {
				mainCmd.Execute()
			})

			assert.Contains(t, out, tt.expect, "Expect: %v\nActual: %v\n", tt.expect, out)
		})
	}
}

func TestPublish(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "https://"+API_HOST+"/namespaces/ns-sampler/gateway", func(r *http.Request) (*http.Response, error) {
		assert.Contains(t, r.URL.Path, "ns-sampler")
		content := map[string]interface{}{
			"configFile": map[string]interface{}{
				"value": configFileContents,
				"options": map[string]interface{}{
					"filename": "config.yaml",
				},
			},
			"dryRun": true,
		}
		jsonBody, _ := json.Marshal(content)
		payload, _ := io.ReadAll(r.Body)
		assert.Equal(t, string(jsonBody), string(payload))

		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"id": "123",
		})
	})

	cwd := t.TempDir()
	ctx := &pkg.AppContext{
		ApiHost:   API_HOST,
		Cwd:       cwd,
		Namespace: "ns-sampler",
	}
	fileName := "config.yaml"
	filePath := filepath.Join(cwd, fileName)
	os.WriteFile(filePath, []byte(configFileContents), 0644)
	opts := &publishOptions{
		configFile: fileName,
		dryRun:     true,
	}
	err := Publish(ctx, opts)
	assert.Nil(t, err, "request success")
}

func TestPublishError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "https://"+API_HOST+"/namespaces/ns-sampler/gateway", func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(401, "Unauthorized")
	})

	cwd := t.TempDir()
	ctx := &pkg.AppContext{
		ApiHost:   API_HOST,
		Cwd:       cwd,
		Namespace: "ns-sampler",
	}
	fileName := "config.yaml"
	filePath := filepath.Join(cwd, fileName)
	os.WriteFile(filePath, []byte(configFileContents), 0644)
	opts := &publishOptions{
		configFile: fileName,
		dryRun:     false,
	}
	err := Publish(ctx, opts)
	assert.ErrorContains(t, err, "Unauthorized")
	assert.NotNil(t, err, "request failed")
}
