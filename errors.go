package interfaces

import (
	"fmt"
	"io"
)

type baseErr struct {
	msg string
}

func NewBasicErr(x string) error {
	return &baseErr{msg: x}
}

func (e *baseErr) Error() string { return e.msg }

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withMessage) Unwrap() error { return w.cause }

// Format provides an implementation of "Formatter", for optimized processing
func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if nil == err {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

func Wrapf(err error, format string, args ...any) error {
	if nil == err {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}
