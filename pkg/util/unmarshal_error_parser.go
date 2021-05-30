package util

import (
	"encoding/json"
	"fmt"
)

func UnmarshalErrorParser(err error) string {
	if typeError, ok := err.(*json.UnmarshalTypeError); ok {
		switch typeError.Type.Name() {
		case "string":
			return fmt.Sprintf("%s must be a string", typeError.Field)
		case "int32", "int64", "float64", "float32":
			return fmt.Sprintf("%s must be a number", typeError.Field)
		}
	}
	return ""
}
