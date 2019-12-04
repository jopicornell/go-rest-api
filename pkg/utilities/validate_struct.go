package utilities

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStruct(ifc interface{}) validator.ValidationErrors {
	if err := validate.Struct(ifc); err != nil {
		return validate.Struct(ifc).(validator.ValidationErrors)
	}
	return nil
}
