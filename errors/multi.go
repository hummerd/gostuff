package errors

import (
	"strings"
)

// NewMultiError returns new MultiError initialized with errs.
func NewMultiError(errs ...error) *MultiError {
	me := &MultiError{
		errs: errs,
	}
	for _, e := range errs {
		if e != nil {
			me.actual++
		}
	}
	return me
}

// MultiError error that combines multiple errors
type MultiError struct {
	errs   []error
	actual int
}

// Error returns combined error message from all errors added to MultiError
func (e *MultiError) Error() string {
	sb := strings.Builder{}
	for _, err := range e.errs {
		if err != nil {
			sb.WriteString(err.Error())
			sb.WriteString("; ")
		}
	}
	return sb.String()
}

// Add - addes error to multi error. If err is nil only AddedLen counter is increased, otherwise
// both - AddedLen and ActualLen are increased.
func (e *MultiError) Add(err error) {
	e.errs = append(e.errs, err)
	if err != nil {
		e.actual++
	}
}

// Get returns error added at i position.
func (e *MultiError) Get(i int) error {
	return e.errs[i]
}

// AddedLen returns number of time add was called.
func (e *MultiError) AddedLen() int {
	return len(e.errs)
}

// ActualLen returns number of actual (non nil) errors added to MultiError.
func (e *MultiError) ActualLen() int {
	return e.actual
}

// IfHasErrors returns e if there were non nil errors added to MultiError. Otherwise returns nil.
func (e *MultiError) IfHasErrors() error {
	if e.actual > 0 {
		return e
	}
	return nil
}
