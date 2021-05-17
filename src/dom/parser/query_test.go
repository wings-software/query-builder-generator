//go:generate goyacc -o query.y.go -p Dom query.y

package parser

import (
	"testing"
  "github.com/stretchr/testify/assert"
)

func TestSanity(t *testing.T) {
	vars := map[string]interface{}{
		"A": 1,
		"B": 1,
	}
	f := Parse("NOT (A IS B)", vars)
	assert.Equal(t, f, false)
}