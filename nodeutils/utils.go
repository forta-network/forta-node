package nodeutils

import "sync"

// ErrorCounter checks incoming errors and tells if we are over
// the max amount of consecutive errors.
type ErrorCounter struct {
	max             uint
	isCriticalError func(error) bool
	count           uint
	sync.Mutex
}

// NewErrorCounter creates a new error counter.
func NewErrorCounter(max uint, isCriticalError func(error) bool) *ErrorCounter {
	return &ErrorCounter{
		max:             max,
		isCriticalError: isCriticalError,
	}
}

func (ec *ErrorCounter) TooManyErrs(err error) bool {
	ec.Lock()
	defer ec.Unlock()
	if err == nil || !ec.isCriticalError(err) {
		ec.count = 0 // reset if other errors or no errors
		return false
	}
	ec.count++
	return ec.count >= ec.max
}
