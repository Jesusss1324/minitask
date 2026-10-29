package services

import (
	"errors"
)

type ValidationError struct {
	Message string
}

func (c *ValidationError) Error() string {
	return c.Message
}

var ErrTaskNotFound = errors.New("task not found")
