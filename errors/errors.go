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

type wrappedError struct {
	msg   string
	cause error
}

// Error returns error's text
func (we *wrappedError) Error() string {
	return we.msg
}

// Cause returns error's cause
func (we *wrappedError) Cause() error {
	return we.cause
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Wrap wraps error with message if err not nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &wrappedError{
		msg:   msg + err.Error(),
		cause: err,
	}
}

// Wrapf wraps error with formatted message if err not nil
func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrappedError{
		msg:   fmt.Sprintf(format, a...) + err.Error(),
		cause: err,
	}
}

// JointError is error returned by Join method in case both errors are not nil.
// Cause for JointError is ErrOne.
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

// Cause returns cause for this error. For JointError cause is ErrOne.
func (e *JointError) Cause() error {
	return e.ErrOne
}

// Join joins two errors in one
// If both errors are nil - returns nil
// If only one error not nil returns that error
// If both not nil returns JointError with both errors (cause for JointError is ErrOne)
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
