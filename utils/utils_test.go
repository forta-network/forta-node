package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapToList(t *testing.T) {
	input := map[string]string{
		"b": "b-value",
		"a": "a-value",
		"c": "c-value",
	}
	list := MapToList(input)
	assert.Len(t, list, 3)

	// should be sorted
	assert.Equal(t, []string{
		"a=a-value", "b=b-value", "c=c-value",
	}, list)
}
