package framework

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator and turns violations into a
// Laravel-style error bag: field name (from the json tag) → human messages.
type Validator struct {
	engine *validator.Validate
}

// NewValidator builds the validation engine. Field names in error bags use
// the struct's json tags so they line up with what the Vue client sent.
func NewValidator() *Validator {
	engine := validator.New(validator.WithRequiredStructEnabled())

	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})

	return &Validator{engine: engine}
}

// Struct validates a DTO and returns the error bag (empty when valid).
func (v *Validator) Struct(dto any) map[string][]string {
	errorBag := map[string][]string{}

	err := v.engine.Struct(dto)
	if err == nil {
		return errorBag
	}

	violations, ok := err.(validator.ValidationErrors)
	if !ok {
		errorBag["_error"] = []string{err.Error()}
		return errorBag
	}

	for _, violation := range violations {
		field := violation.Field()
		errorBag[field] = append(errorBag[field], message(violation))
	}
	return errorBag
}

// message renders a readable message for a single violation.
func message(v validator.FieldError) string {
	label := humanize(v.Field())

	switch v.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", label)
	case "email":
		return fmt.Sprintf("The %s must be a valid email address.", label)
	case "min":
		if v.Kind() == reflect.String {
			return fmt.Sprintf("The %s must be at least %s characters.", label, v.Param())
		}
		return fmt.Sprintf("The %s must be at least %s.", label, v.Param())
	case "max":
		if v.Kind() == reflect.String {
			return fmt.Sprintf("The %s may not be greater than %s characters.", label, v.Param())
		}
		return fmt.Sprintf("The %s may not be greater than %s.", label, v.Param())
	case "eqfield":
		return fmt.Sprintf("The %s must match the %s field.", label, humanize(v.Param()))
	case "url":
		return fmt.Sprintf("The %s must be a valid URL.", label)
	case "numeric":
		return fmt.Sprintf("The %s must be a number.", label)
	case "oneof":
		return fmt.Sprintf("The selected %s is invalid.", label)
	default:
		return fmt.Sprintf("The %s field is invalid (%s).", label, v.Tag())
	}
}

// humanize turns "passwordConfirmation" / "password_confirmation" into
// "password confirmation".
func humanize(field string) string {
	var b strings.Builder
	for i, r := range field {
		switch {
		case r == '_' || r == '-':
			b.WriteRune(' ')
		case r >= 'A' && r <= 'Z':
			if i > 0 {
				b.WriteRune(' ')
			}
			b.WriteRune(r - 'A' + 'a')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
