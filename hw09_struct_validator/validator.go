package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type NoStructTypeError struct {
	message string
}

func (n NoStructTypeError) Error() string {
	return n.message
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidatorError struct {
	Field string
	Err   error
}

type (
	ValidationErrors  []ValidationError
	ValidatorSetError []ValidatorError
)

var (
	ErrValidate = errors.New("validate error")
	ErrTags     = errors.New("program fail")
)

func (v ValidationErrors) Error() string {
	var errorOut string
	for _, e := range v {
		errorOut += fmt.Sprintf("err in field - '%v': %s\n", e.Field, e.Err)
	}
	return errorOut
}

func (v ValidatorSetError) Error() string {
	var errorOut string
	for _, e := range v {
		errorOut += fmt.Sprintf("err in field - '%v': %s\n", e.Field, e.Err)
	}
	return errorOut
}

func valueIn(val any, set string) bool {
	setList := strings.Split(set, ",")
	switch val.(type) {
	case int:
		for _, s := range setList {
			num, _ := strconv.Atoi(s)
			if val == num {
				return true
			}
		}
		return false
	case string:
		for _, str := range setList {
			if val == str {
				return true
			}
		}
		return false
	}
	return false
}

func validateString(commands []string, str string) error {
	switch commands[0] {
	case "in":
		if valueIn(str, commands[1]) {
			return nil
		}
		return fmt.Errorf("%w, value %v not in %v", ErrValidate, str, commands[1])
	case "len":
		requireLen, err := strconv.Atoi(commands[1])
		if err != nil {
			return err
		}
		if len(str) == requireLen {
			return nil
		}
		return fmt.Errorf("%w, length of the %v is not equal %v", ErrValidate, str, commands[1])
	case "regexp":
		exp, err := regexp.Compile(commands[1])
		if err != nil {
			return err
		}
		if exp.Match([]byte(str)) {
			return nil
		}
		return fmt.Errorf("%w, %s is not match for %v expression", ErrValidate, str, commands[1])
	}
	return fmt.Errorf("%w, %s validator is not supported for type string", ErrTags, commands[0])
}

func validateInt(commands []string, number int) error {
	switch commands[0] {
	case "in":
		if valueIn(number, commands[1]) {
			return nil
		}
		return fmt.Errorf("%w, value %v not in %v", ErrValidate, number, commands[1])
	case "max":
		max, err := strconv.Atoi(commands[1])
		if err != nil {
			return err
		}
		if number > max {
			return fmt.Errorf("%w, number %v over that %q", ErrValidate, number, commands[1])
		}
		return nil
	case "min":
		min, err := strconv.Atoi(commands[1])
		if err != nil {
			return err
		}
		if number < min {
			return fmt.Errorf("%w, number %v less that %v", ErrValidate, number, commands[1])
		}
		return nil
	}
	return fmt.Errorf("%w, %s validator is not supported for type int", ErrTags, commands[0])
}

func validateChecks(field *reflect.StructField,
	fieldValue reflect.Value,
	rules []string,
) (ValidationErrors, ValidatorSetError) {
	var validErrors ValidationErrors
	var validSetError ValidatorSetError

	for _, item := range rules {
		var err error
		commands := strings.Split(item, ":")
		switch fieldValue.Kind().String() {
		case "string":
			err = validateString(commands, fieldValue.String())
		case "int":
			err = validateInt(commands, int(fieldValue.Int()))
		}
		if err != nil {
			if errors.Is(err, ErrTags) {
				validSetError = append(validSetError, ValidatorError{Field: field.Name, Err: err})
			} else if errors.Is(err, ErrValidate) {
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}
	return validErrors, validSetError
}

func parseLabel(labelString string) []string {
	labelList := strings.Split(labelString, "|")
	return labelList
}

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind().String() != "struct" {
		return NoStructTypeError{message: fmt.Sprintf("type error: given object of type '%v', expected %v",
			reflect.TypeOf(v).Kind().String(), "struct")}
	}

	var ResultErrorsList ValidationErrors
	var ValidationsErrTags ValidatorSetError

	valueFields := reflect.ValueOf(v)

	for i := 0; i < valueFields.Type().NumField(); i++ {
		field := reflect.TypeOf(v).Field(i)
		fieldValue := valueFields.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		if fieldValue.Kind() == reflect.Slice {
			for j := 0; j < valueFields.Field(i).Len(); j++ {
				REL, VTE := validateChecks(&field, fieldValue.Index(j), parseLabel(tag))
				ResultErrorsList = append(ResultErrorsList, REL...)
				ValidationsErrTags = append(ValidationsErrTags, VTE...)
			}
			continue
		}
		REL, VTE := validateChecks(&field, fieldValue, parseLabel(tag))
		ResultErrorsList = append(ResultErrorsList, REL...)
		ValidationsErrTags = append(ValidationsErrTags, VTE...)
	}

	if len(ValidationsErrTags) != 0 {
		return ValidationsErrTags
	} else if len(ResultErrorsList) != 0 {
		return ResultErrorsList
	}
	return nil
}
