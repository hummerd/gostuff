package errors

import (
	"strings"
	"sync"
)

// SyncMultiError is syncronized version of MultiError
type SyncMultiError struct {
	sync.RWMutex
	errs   []error
	actual int
}

// Error returns combined error message from all errors added to SyncMultiError
func (e *SyncMultiError) Error() string {
	e.RLock()
	defer e.RUnlock()

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
func (e *SyncMultiError) Add(err error) {
	e.Lock()
	defer e.Unlock()

	e.errs = append(e.errs, err)
	if err != nil {
		e.actual++
	}
}

// Get returns error added at i position.
func (e *SyncMultiError) Get(i int) error {
	e.RLock()
	defer e.RUnlock()

	return e.errs[i]
}

// AddedLen returns number of time add was called.
func (e *SyncMultiError) AddedLen() int {
	e.RLock()
	defer e.RUnlock()

	return len(e.errs)
}

// ActualLen returns number of actual (non nil) errors added to MultiError.
func (e *SyncMultiError) ActualLen() int {
	e.RLock()
	defer e.RUnlock()

	return e.actual
}

// IfHasErrors returns e if there were non nil errors added to SyncMultiError. Otherwise returns nil.
func (e *SyncMultiError) IfHasErrors() error {
	e.RLock()
	defer e.RUnlock()

	if e.actual > 0 {
		return e
	}
	return nil
}
