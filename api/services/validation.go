package services

import (
	"strings"
	"unicode/utf8"
)

func validateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)

	length := utf8.RuneCountInString(title)

	if length < 3 || length > 100 {
		return "", &ValidationError{
			Message: "El título debe tener entre 3 y 100 caracteres.",
		}
	}

	return title, nil
}

func validateDescription(description string) (string, error) {
	length := utf8.RuneCountInString(description)

	if length > 1000 {
		return "", &ValidationError{
			Message: "La descripción no puede superar los 1000 caracteres.",
		}
	}

	return description, nil
}

func validateEnum(value string, allowed []string, message string) error {
	for _, option := range allowed {
		if value == option {
			return nil
		}
	}

	return &ValidationError{Message: message}
}
