package utils

import (
	"sync"

	"github.com/tylertreat/BoomFilters"
)

type Cache interface {
	Add(str string)
	Exists(str string) bool
	ExistsAndAdd(str string) bool
}

type cache struct {
	filter boom.Filter
	mux    *sync.Mutex
}

func (c *cache) Add(str string) {
	c.filter.Add([]byte(str))
}

//ExistsAndAdd returns true if already exists, otherwise adds the item to cache and returns false
func (c *cache) ExistsAndAdd(str string) bool {
	return c.filter.TestAndAdd([]byte(str))
}

//Exists returns true if already exists
func (c *cache) Exists(str string) bool {
	return c.filter.Test([]byte(str))
}

//NewCache creates a new cache
func NewCache(size uint) *cache {
	bf := boom.NewInverseBloomFilter(size)
	return &cache{
		filter: bf,
		mux:    &sync.Mutex{},
	}
}
