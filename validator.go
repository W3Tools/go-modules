package gm

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(v any) error {
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var errString string
		for _, err1 := range err.(validator.ValidationErrors) {
			tmp := fmt.Sprintf(
				"%s %s %s [%v]",
				err1.StructField(), err1.Tag(), err1.Param(), err1.Value())
			if len(errString) > 0 {
				errString += " | "
			}
			errString += tmp
		}
		err = fmt.Errorf(errString)
	}
	return err
}

func ValidateJsonUnmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	err = ValidateStruct(v)
	if err != nil {
		return err
	}
	return nil
}
