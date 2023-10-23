package cmd

import (
	"fmt"
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
		namespace  string
	}{
		{
			name:       "successful straight publish",
			setup:      nil,
			configFile: "config.yaml",
			response:   httpmock.NewStringResponder(200, `{"id": 1}`),
			expect:     "Gateway config published",
			args:       []string{"config.yaml"},
			namespace:  "ns-sampler",
		},
		{
			name:       "api error",
			setup:      nil,
			configFile: "config.yaml",
			response:   httpmock.NewStringResponder(500, `{"error": "something went wrong"}`),
			expect:     "something went wrong",
			args:       []string{"config.yaml"},
			namespace:  "ns-sampler",
		},
		{
			name:       "missing namespace",
			setup:      nil,
			configFile: "config.yaml",
			response:   httpmock.NewStringResponder(500, `{"error": "something went wrong"}`),
			expect:     "No namespace has been set",
			args:       []string{"config.yaml"},
			namespace:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("PUT", "https://"+API_HOST+"/gw/api/namespaces/ns-sampler/gateway", tt.response)
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
				Namespace: tt.namespace,
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

func TestPrepareConfigFile(t *testing.T) {
	cwd := t.TempDir()
	ctx := &pkg.AppContext{
		ApiHost:   API_HOST,
		Cwd:       cwd,
		Namespace: "ns-sampler",
	}
	fileName := "config.yaml"
	filePath := filepath.Join(cwd, fileName)
	os.WriteFile(filePath, []byte(configFileContents), 0644)
	opts := &PublishGatewayOptions{
		inputs: []string{fileName},
		dryRun: true,
	}
	config, err := PrepareConfigFile(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	actualBytes, err := io.ReadAll(config)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes)
	assert.Equal(t, configFileContents, actual, "it returns a config file")
}

func TestMultiPrepareConfigFile(t *testing.T) {
	cwd := t.TempDir()
	ctx := &pkg.AppContext{
		Cwd: cwd,
	}
	for i := range "123" {
		fileName := fmt.Sprintf("config-%d.yaml", i)
		contents := fmt.Sprintf(`
_format_version: "1.1"
services:
  - name: Demo_App_%d
    url: /api/demoapp-%d
    plugins: []`, i, i)
		filePath := filepath.Join(cwd, fileName)
		os.WriteFile(filePath, []byte(contents), 0755)
	}
	opts := &PublishGatewayOptions{
		inputs: []string{""},
		dryRun: false,
	}
	config, err := PrepareConfigFile(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	actualBytes, err := io.ReadAll(config)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes)
	combined := []byte(`
_format_version: "1.1"
services:
  - name: Demo_App_0
    url: /api/demoapp-0
    plugins: []
---

_format_version: "1.1"
services:
  - name: Demo_App_1
    url: /api/demoapp-1
    plugins: []
---

_format_version: "1.1"
services:
  - name: Demo_App_2
    url: /api/demoapp-2
    plugins: []`)
	expected := string(combined)
	assert.Equal(t, expected, actual, "it returns a multi-document yaml file")
}

func TestMultPrepareEmptyDir(t *testing.T) {
	cwd := t.TempDir()
	ctx := &pkg.AppContext{
		Cwd: cwd,
	}
	opts := &PublishGatewayOptions{
		inputs: []string{cwd},
	}

	_, err := PrepareConfigFile(ctx, opts)
	assert.Error(t, err, "There is no yaml files in this directory")
}

func TestIncorrectFileType(t *testing.T) {
	opts := &PublishGatewayOptions{
		inputs: []string{"test.json"},
	}

	_, err := PrepareConfigFile(ctx, opts)
	assert.Error(t, err, "non-yaml file types not allowed")
}

func TestMixedFileArguments(t *testing.T) {
	cwd := t.TempDir()
	os.Mkdir(filepath.Join(cwd, "other-configs"), 0755)
	os.Mkdir(filepath.Join(cwd, "nested"), 0755)
	err := os.WriteFile(filepath.Join(cwd, "config.yaml"), []byte("config: 1"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(cwd, "other-configs", "config.yaml"), []byte("config: 2"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(cwd, "nested", "config.yaml"), []byte("config: 3"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	ctx := &pkg.AppContext{
		Cwd: cwd,
	}

	opts := &PublishGatewayOptions{
		inputs: []string{"config.yaml", "other-configs/", "nested/config.yaml"},
	}
	config, err := PrepareConfigFile(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	actualBytes, err := io.ReadAll(config)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes)
	combined := []byte(`config: 1
---
config: 2
---
config: 3`)
	expected := string(combined)
	assert.Equal(t, expected, actual, "it crawls several entry points")
}

func TestPublishGatewayWithQualifier(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "https://"+API_HOST+"/gw/api/namespaces/ns-sampler/gateway", func(r *http.Request) (*http.Response, error) {
		assert.Contains(t, r.URL.Path, "ns-sampler")
		assert.Equal(t, "myqualifier", r.FormValue("qualifier"))

		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"message": "gateway published",
			"results": "aok",
			"error":   "",
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
	opts := &PublishGatewayOptions{
		inputs:    []string{fileName},
		qualifier: "myqualifier",
		dryRun:    true,
	}
	_, err := PrepareConfigFile(ctx, opts)
	assert.Nil(t, err, "request success")
}
