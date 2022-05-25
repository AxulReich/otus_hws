package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateKeyWord         = "validate"
	tagsDelimiter           = " "
	tagNameSep              = ":"
	ruleOrSep               = "|"
	validateRuleParentheses = `"`

	stringRuleLen = "len"
	stringRuleReg = "regexp"
	stringRuleIn  = "in"
	intRuleMax    = "max"
	intRuleMin    = "min"
	intRuleIn     = stringRuleIn
)

type (
	ValidationError struct {
		Field string
		Err   error
	}

	ValidationErrors []ValidationError
	validateRules    = map[string]string
)

var (
	ErrNotStruct       = errors.New("not struct passed")
	ErrUnexpectedRule  = errors.New("unexpected rule")
	ErrCastRuleValue   = errors.New("can't cast rule value")
	ErrUnsupportedType = errors.New("type is unsupported")

	ErrInvalidStringLength    = errors.New("string field value has invalid length")
	ErrInvalidStringSignature = errors.New("string field value doesn't match regular expression")
	ErrInvalidStringValue     = errors.New("string field value doesn't in set")
	ErrInvalidIntValue        = errors.New("int field value doesn't in set")
	ErrInvalidIntMin          = errors.New("int field value is less than min")
	ErrInvalidIntMax          = errors.New("int field value is greater than max")
)

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
		return ErrNotStruct
	}
	var validationErrors ValidationErrors
	for i := 0; i < rv.NumField(); i++ {
		// if rv.Field(i).CanInterface() == false -> field is unexported
		if !rv.Field(i).CanInterface() {
			continue
		}

		tags := string(rv.Type().Field(i).Tag)

		for _, tag := range strings.Split(tags, tagsDelimiter) {
			if strings.Split(tag, tagNameSep)[0] == validateKeyWord {
				trimmedTag := strings.Trim(tag[len(validateKeyWord+tagNameSep):], validateRuleParentheses)
				field := rv.Field(i)
				fieldName := rv.Type().Field(i).Name
				err := validateField(trimmedTag, field)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   err,
					})
				}
			}
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func validateField(validationRules string, val reflect.Value) error {
	var (
		rulesRaw = strings.Split(validationRules, ruleOrSep)
		rules    = make(validateRules)
	)

	for _, tag := range rulesRaw {
		ruleName := strings.Split(tag, tagNameSep)[0]
		ruleValue := strings.Split(tag, tagNameSep)[1]
		rules[ruleName] = ruleValue
	}
	// nolint:exhaustive
	switch val.Kind() {
	case reflect.Int:
		err := validateIntField(rules, []int64{val.Int()})
		if err != nil {
			return err
		}
		return nil
	case reflect.String:
		err := validateStringField(rules, []string{val.String()})
		if err != nil {
			return err
		}
		return nil
	case reflect.SliceOf(val.Type()).Kind():
		switch t := val.Interface().(type) {
		case []string:
			var value []string
			value = append(value, t...)
			if err := validateStringField(rules, value); err != nil {
				return err
			}
			return nil
		case []int:
			var value []int64
			for _, v := range t {
				value = append(value, int64(v))
			}
			if err := validateIntField(rules, value); err != nil {
				return err
			}
		default:
			return ErrUnsupportedType
		}
	default:
		return ErrUnsupportedType
	}
	return nil
}

func validateStringField(rules validateRules, value []string) error {
	for _, v := range value {
		for ruleName, ruleValue := range rules {
			switch ruleName {
			case stringRuleIn:
				if !strings.Contains(ruleValue, v) {
					return ErrInvalidStringValue
				}
			case stringRuleLen:
				mustLen, err := strconv.Atoi(ruleValue)
				if err != nil {
					return ErrCastRuleValue
				}
				if len([]rune(v)) != mustLen {
					return ErrInvalidStringLength
				}
			case stringRuleReg:
				ruleValue = strings.ReplaceAll(ruleValue, "\\\\", `\`)
				reg := regexp.MustCompile(ruleValue)
				if !reg.MatchString(v) {
					return ErrInvalidStringSignature
				}
			default:
				return ErrUnexpectedRule
			}
		}
	}
	return nil
}

func validateIntField(rules validateRules, value []int64) error {
	for _, v := range value {
		for ruleName, ruleValue := range rules {
			switch ruleName {
			case intRuleMax:
				max, err := strconv.ParseInt(ruleValue, 10, 64)
				if err != nil {
					return ErrCastRuleValue
				}
				if v > max {
					return ErrInvalidIntMax
				}
			case intRuleMin:
				min, err := strconv.ParseInt(ruleValue, 10, 64)
				if err != nil {
					return ErrCastRuleValue
				}
				if v < min {
					return ErrInvalidIntMin
				}
			case intRuleIn:
				if !strings.Contains(ruleValue, strconv.FormatInt(v, 10)) {
					return ErrInvalidIntValue
				}
			default:
				return ErrUnexpectedRule
			}
		}
	}
	return nil
}
