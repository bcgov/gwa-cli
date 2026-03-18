package cmd

import (
	"testing"

	"github.com/bcgov/gwa-cli/pkg"
	"github.com/stretchr/testify/assert"
)

func TestIsLocalHost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		apiHost  string
		expected bool
	}{
		{
			name:     "oauth2proxy.localtest.me with port",
			host:     "",
			apiHost:  "oauth2proxy.localtest.me:4180",
			expected: true,
		},
		{
			name:     "localtest.me subdomain",
			host:     "app.localtest.me",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "localhost",
			host:     "localhost",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "localhost with port",
			host:     "localhost:8080",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "127.0.0.1",
			host:     "127.0.0.1",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "127.0.0.1 with port",
			host:     "127.0.0.1:3000",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "case insensitive localtest.me",
			host:     "LOCALTEST.ME",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "case insensitive localhost",
			host:     "LocalHost",
			apiHost:  "",
			expected: true,
		},
		{
			name:     "gov.bc.ca dev api - not local",
			host:     "",
			apiHost:  "api-gov-bc-ca.dev.api.gov.bc.ca",
			expected: false,
		},
		{
			name:     "gov.bc.ca prod api - not local",
			host:     "",
			apiHost:  "api.gov.bc.ca",
			expected: false,
		},
		{
			name:     "aps.gov.bc.ca - not local",
			host:     "",
			apiHost:  "aps.gov.bc.ca",
			expected: false,
		},
		{
			name:     "empty host - not local",
			host:     "",
			apiHost:  "",
			expected: false,
		},
		{
			name:     "Host overrides ApiHost when both set",
			host:     "oauth2proxy.localtest.me:4180",
			apiHost:  "api-gov-bc-ca.dev.api.gov.bc.ca",
			expected: true,
		},
		{
			name:     "Host overrides ApiHost - remote wins",
			host:     "api-gov-bc-ca.dev.api.gov.bc.ca",
			apiHost:  "localhost",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &pkg.AppContext{
				Host:    tt.host,
				ApiHost: tt.apiHost,
			}
			result := isLocalHost(ctx)
			assert.Equal(t, tt.expected, result, "isLocalHost(ctx) = %v, want %v", result, tt.expected)
		})
	}
}
