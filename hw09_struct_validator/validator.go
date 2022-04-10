package hw09structvalidator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	s := strings.Builder{}

	for _, e := range v {
		s.WriteString(fmt.Sprintf("field: %q - %s\n", e.Field, e.Err.Error()))
	}

	return s.String()
}

func Validate(v interface{}) error {
	if ref
	switch v.(type) {
	case :

	}
	// Place your code here.
	return nil
}
