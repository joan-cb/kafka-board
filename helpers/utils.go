package helpers

import (
	"kafka-board/types"

	"github.com/xeipuuv/gojsonschema"
)

// CheckErr is a helper function to check if an error is present

func CheckErr(e error) bool {
	return e != nil
}

func ValidatePayload(payload interface{}, schema types.Schema) (bool, []string, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema.Schema)
	documentLoader := gojsonschema.NewGoLoader(payload)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, nil, err
	}

	if !result.Valid() {
		var errors []string
		for _, err := range result.Errors() {
			errors = append(errors, err.String())
		}
		return false, errors, nil
	}

	return true, nil, nil
}
