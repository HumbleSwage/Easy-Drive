package types

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitValidate() {
	Validate = validator.New()
}

func GetValidate() *validator.Validate {
	return Validate
}
