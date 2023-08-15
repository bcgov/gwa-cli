package pkg

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKebabCase(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "spaces to kebab",
			input:  "my service   name",
			expect: "my-service-name",
		},
		{
			name:   "snake to kebab",
			input:  "my_service_name",
			expect: "my-service-name",
		},
		{
			name:   "weird chars to kebab",
			input:  "my*service^name",
			expect: "my-service-name",
		},
		{
			name:   "kebab spaces",
			input:  "MY-SERVICE-NAME",
			expect: "my-service-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, KebabCase(tt.input))
		})
	}
}

func TestStartCase(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "all lower case",
			input:  "my service name",
			expect: "My service name",
		},
		{
			name:   "all uppercase case",
			input:  "MY SERVICE NAME",
			expect: "My service name",
		},
		{
			name:   "when snake case",
			input:  "my-service-name",
			expect: "My-service-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, StartCase(tt.input))
		})
	}
}

func TestAppId(t *testing.T) {
	assert.Regexp(t, regexp.MustCompile(`[A-Z0-9]{12}`), AppId(12))
	assert.Regexp(t, regexp.MustCompile(`[A-Z0-9]{6}`), AppId(6))
}
