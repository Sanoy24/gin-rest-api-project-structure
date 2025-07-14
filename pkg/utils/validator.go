package utils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

func FormatValidationError(err error, obj any) map[string]string {
	errors := make(map[string]string)

	objType := reflect.TypeOf(obj)
	fmt.Println(objType, "|", objType.Kind(), "|", reflect.Ptr, "|", objType.Elem())

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		objType := reflect.TypeOf(obj)
		fmt.Println(objType, "|", objType.Kind(), "|", reflect.Ptr, "|", objType.Elem())
		if objType.Kind() == reflect.Ptr {
			objType = objType.Elem()
		}

		for _, e := range validationErrors {
			// Get the struct field by name
			if field, found := objType.FieldByName(e.StructField()); found {
				// Get the json tag
				jsonKey := field.Tag.Get("json")
				if jsonKey == "" || jsonKey == "-" {
					jsonKey = e.Field()
				}
				errors[jsonKey] = getErrorMessage(e)
			} else {
				errors[e.Field()] = getErrorMessage(e)
			}
		}
	}
	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "value is too long"
	case "oneof":
		return "Invalid value"
	default:
		return "Invalid value"
	}
}
