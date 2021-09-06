package agentpool

import "sync"

// errorCounter checks incoming errors and tells if we are over
// the max amount of consecutive errors.
type errorCounter struct {
	max      uint
	errCheck func(error) bool
	count    uint
	sync.Mutex
}

// NewErrorCounter creates a new error counter.
func NewErrorCounter(max uint, errCheck func(error) bool) *errorCounter {
	return &errorCounter{
		max:      max,
		errCheck: errCheck,
	}
}

func (ec *errorCounter) TooManyErrs(err error) bool {
	ec.Lock()
	defer ec.Unlock()
	if err == nil || !ec.errCheck(err) {
		ec.count = 0 // reset if other errors or no errors
		return false
	}
	ec.count++
	return ec.count >= ec.max
}
