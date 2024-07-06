package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ResultErrorsList ValidationErrors

type ValidationError struct {
	Field string
	Err   error
}

type NoStructTypeError struct {
	message string
}

type (
	ValidationErrors []ValidationError
	NoStructType     error
)

func (n NoStructTypeError) Error() string {
	panic(n.message)
}

func (v ValidationErrors) Error() string {
	var panicOut string
	for _, e := range v {
		panicOut = panicOut + fmt.Sprintf("error in '%v' field: %s\n", e.Field, e.Err)
	}
	panic(panicOut)
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

func validateChecks(field *reflect.StructField, fieldValue reflect.Value, rules []string) ValidationErrors {
	var validErrors ValidationErrors

	for _, item := range rules {
		commands := strings.Split(item, ":")
		if fieldValue.Kind().String() == "string" {
			switch commands[0] {
			case "in":
				if valueIn(fieldValue.String(), commands[1]) {
					continue
				}
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("value %v not in %v", fieldValue.String(), commands[1])})
			case "len":
				requireLen, err := strconv.Atoi(commands[1])
				if err != nil {
					validErrors = append(validErrors, ValidationError{Field: field.Name, Err: err})
					continue
				}
				if len(fieldValue.String()) == requireLen {
					continue
				}
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("length of the %s is not equal %v", fieldValue.String(), commands[1])})
			case "regexp":
				exp, err := regexp.Compile(commands[1])
				if err != nil {
					validErrors = append(validErrors, ValidationError{Field: field.Name, Err: err})
				}
				if exp.Match([]byte(fieldValue.String())) {
					continue
				}
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("%s is not match for %v expression", fieldValue.String(), commands[1])})
			default:
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("program fail: %s validator is not supported for type string", commands[0])})
			}
		} else if fieldValue.Kind().String() == "int" {
			switch commands[0] {
			case "in":
				if valueIn(int(fieldValue.Int()), commands[1]) {
					continue
				}
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("value %v not in %q", fieldValue.Int(), commands[1])})
			case "min":
				min, err := strconv.Atoi(commands[1])
				if err != nil {
					validErrors = append(validErrors, ValidationError{Field: field.Name, Err: err})
					continue
				}
				if int(fieldValue.Int()) < min {
					validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("number %v less that %v", fieldValue.Int(), commands[1])})
				}
			case "max":
				max, err := strconv.Atoi(commands[1])
				if err != nil {
					validErrors = append(validErrors, ValidationError{Field: field.Name, Err: err})
					continue
				}
				if int(fieldValue.Int()) > max {
					validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("number %v over that %q", fieldValue.Int(), commands[1])})
				}
			default:
				validErrors = append(validErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("program fail: %s validator is not supported for type int", commands[0])})
			}
		}
	}
	return validErrors
}

func parseLabel(labelString string) []string {
	labelList := strings.Split(labelString, "|")
	return labelList
}

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind().String() != "struct" {
		return NoStructTypeError{message: fmt.Sprintf("type error: given object of type '%v', expected %v", reflect.TypeOf(v).Kind().String(), "struct")}
	}

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
				ResultErrorsList = append(ResultErrorsList, validateChecks(&field, fieldValue.Index(j), parseLabel(tag))...)
			}
			continue
		}
		ResultErrorsList = append(ResultErrorsList, validateChecks(&field, fieldValue, parseLabel(tag))...)
	}

	if len(ResultErrorsList) == 0 {
		return nil
	}
	return ResultErrorsList
}
