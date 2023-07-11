package validator

import (
	validator "github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func ProvideValidator() fx.Option {
	return fx.Provide(NewValidator)
}

func NewValidator() *validator.Validate {
	return validator.New()
}
