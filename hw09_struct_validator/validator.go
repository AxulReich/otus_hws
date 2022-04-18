package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

// var NotStructError = errors.New("not struct passed")

func (v ValidationErrors) Error() string {
	s := strings.Builder{}

	for _, e := range v {
		s.WriteString(fmt.Sprintf("field: %q - %s\n", e.Field, e.Err.Error()))
	}

	return s.String()
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expect struct, but recieved: %T", v)
	}
	// Place your code here.
	return nil
}
