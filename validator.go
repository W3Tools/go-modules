package gm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Verify the structure
func ValidateStruct(v interface{}) error {
	validate := validator.New()

	if err := validate.Struct(v); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var msg strings.Builder
		for _, err1 := range err.(validator.ValidationErrors) {
			if msg.Len() > 0 {
				msg.WriteString(" | ")
			}
			msg.WriteString(fmt.Sprintf("%s %s %s [%v]", err1.StructField(), err1.Tag(), err1.Param(), err1.Value()))
		}
		return fmt.Errorf(msg.String())
	}
	return nil
}

// Validate JSON data and parse it into the specified structure
func ValidateJsonUnmarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return ValidateStruct(v)
}
