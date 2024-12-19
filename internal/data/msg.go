package data

import "errors"

type ErrMsg struct{ err error }

func (e ErrMsg) Error() string { return e.err.Error() }

func WrapError(err error) ErrMsg {
	return ErrMsg{err}
}

func NewError(s string) ErrMsg {
	return WrapError(errors.New(s))
}
