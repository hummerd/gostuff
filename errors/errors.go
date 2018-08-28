package errors

import (
	std_errors "errors"
	"fmt"
)

// New creates new error
func New(msg string) error {
	return std_errors.New(msg)
}

// Newf creates new error with formatted message
func Newf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// Wrap wraps error with message if err not nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return std_errors.New(msg + err.Error())
}

// Wrapf wraps error with formatted message if err not nil
func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return std_errors.New(fmt.Sprintf(format, a...) + err.Error())
}

// JointError error returned by Join method in case both errors are not nil
type JointError struct {
	ErrOne error
	ErrTwo error
}

// Error return error descripton for errOne, errTwo or both of them
func (e *JointError) Error() string {
	if e.ErrTwo != nil {
		if e.ErrOne != nil {
			return e.ErrOne.Error() + "; also: " + e.ErrTwo.Error()
		}
		return e.ErrTwo.Error()
	}
	return ""
}

// Join joins two errors in one
// If both errors are nil - returns nil
// If only one error not nil returns that error
// If both not nil returns JointError with both errors
// For convinient use in defer func () { err = errors.Join(err, some.Close()) }()
func Join(err1, err2 error) error {
	if err1 == nil && err2 == nil {
		return nil
	}
	if err1 != nil && err2 == nil {
		return err1
	}
	if err2 != nil && err1 == nil {
		return err2
	}
	return &JointError{err1, err2}
}
