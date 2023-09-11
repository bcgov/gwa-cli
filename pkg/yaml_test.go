package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var multiDocFile = `key: value
list:
    - one
    - two
---
key: another value
---
lastKey: 123`

func TestSplitYaml(t *testing.T) {
	files, err := SplitYAML([]byte(multiDocFile))
	if err != nil {
		t.Fatal(err)
	}
	var expect [][]byte
	expect = append(expect, []byte(`key: value
list:
    - one
    - two
`))
	expect = append(expect, []byte(`key: another value
`))
	expect = append(expect, []byte(`lastKey: 123
`))

	assert.Equal(t, expect, files)
}

func TestSingleYaml(t *testing.T) {
	var config = []byte(`key: value
list:
    - one
    - two
`)
	files, err := SplitYAML(config)
	if err != nil {
		t.Fatal(err)
	}
	var expect [][]byte
	expect = append(expect, config)
	assert.Equal(t, expect, files)
}
