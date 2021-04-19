package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_Add(t *testing.T) {
	key := "0xe007a511ca2727b330e92a9609ca723868284dcf2d0c0c3009e9c0a4381144a8"
	cache := NewCache(10)
	assert.False(t, cache.Exists(key), "entry should exist")
	cache.Add(key)
	assert.True(t, cache.Exists(key), "entry should exist")
}
