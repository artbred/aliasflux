package internal

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"reflect"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			alias := err.Field()

			field, ok := reflect.TypeOf(i).Elem().FieldByName(alias)
			if ok {
				fieldName := field.Tag.Get("alias")
				if fieldName != "" {
					alias = fieldName
				}
			}

			switch err.Tag() {
			case "required":
				return errors.New(fmt.Sprintf("%s is required", alias))
			case "max":
				return errors.New(fmt.Sprintf("%s is too long", alias))
			case "uuid":
				return errors.New(fmt.Sprintf("%s is not a valid uuid", alias))
			case "email":
				return errors.New("You provided invalid email")
			default:
				return errors.New(fmt.Sprintf("Value for %s must be %s", alias, err.Tag()))
			}
		}
	}

	return nil
}

func SetupValidator(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
}
