package workers_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/forta-protocol/forta-node/utils/workers"
	"github.com/stretchr/testify/require"
)

func TestWorkers(t *testing.T) {
	r := require.New(t)

	w := workers.New(1)

	var (
		testValues1 = []interface{}{"foo"}
		testErr1    error

		testValues2 = []interface{}{}
		testErr2    = errors.New("bar")
	)

	var wg sync.WaitGroup
	wg.Add(2)

	var (
		values1 []interface{}
		err1    error
	)
	go func() {
		output := w.Execute(func() ([]interface{}, error) {
			return testValues1, testErr1
		})
		values1 = output.Values
		err1 = output.Error
		wg.Done()
	}()

	var (
		values2 []interface{}
		err2    error
	)
	go func() {
		output := w.Execute(func() ([]interface{}, error) {
			return testValues2, testErr2
		})
		values2 = output.Values
		err2 = output.Error
		wg.Done()
	}()

	wg.Wait()

	r.Equal(testValues1, values1)
	r.Equal(testErr1, err1)

	r.Equal(testValues2, values2)
	r.Equal(testErr2, err2)
}
