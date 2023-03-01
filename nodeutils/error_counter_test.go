package nodeutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorCounter(t *testing.T) {
	r := require.New(t)

	// max two consecutive non-nil errors means a red flag
	counter := NewErrorCounter(2, func(err error) bool {
		return err != nil
	})

	r.False(counter.TooManyErrs(errors.New("error 1")))
	r.False(counter.TooManyErrs(nil)) // no error - back to zero
	r.False(counter.TooManyErrs(errors.New("error 1")))
	r.True(counter.TooManyErrs(errors.New("error 2"))) // hit the limit - we're done
}
