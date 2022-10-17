package config

import "fmt"

type Error struct {
	reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("unable to build configuration object: %s", e.reason)
}

func NewError(reason string) Error {
	return Error{
		reason: reason,
	}
}
