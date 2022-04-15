package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func Decode(input interface{}, output interface{}) error {
	err := mapstructure.Decode(input, output)
	if err != nil {
		return err
	}
	err = validator.New().Struct(output)
	return err
}
