package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentsSliceToString(t *testing.T) {
	assert.Equal(t, "datasets and issuers", ArgumentsSliceToString([]string{"datasets", "issuers"}, "and"), "Only 2 items results in and as the separator")
	assert.Equal(t, "datasets, issuers, organizations and products", ArgumentsSliceToString([]string{"datasets", "issuers", "organizations", "products"}, "and"), "More than 2 items uses commas until the last separator")
	assert.Equal(t, "datasets or issuers", ArgumentsSliceToString([]string{"datasets", "issuers"}, "or"), "Optional separator")
}
