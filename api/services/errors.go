package services

type ValidationError struct {
	Message string
}

func (c *ValidationError) Error() string {
	return c.Message
}
