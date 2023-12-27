package helper

import (
	"reflect"
)

type Response struct {
	Status bool        `json:"status"`
	Errors interface{} `json:"errors"`
	Data   interface{} `json:"data"`
}

func BuildResponse(data interface{}) Response {
	res := Response{
		Status: true,
		Errors: nil,
		Data:   data,
	}
	return res
}

func BuildErrorResponse(err interface{}) Response {
	res := Response{
		Status: false,
		Errors: err,
		Data:   nil,
	}
	return res
}

func BuildUpdateData(data interface{}) map[string]interface{} {
	updateData := make(map[string]interface{})
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		tag := typeField.Tag.Get("gorm")

		// Skip if field is not exported or has no gorm tag
		if !field.CanInterface() || tag == "-" {
			continue
		}

		// Get field name from gorm tag, or use struct field name
		var name string
		if tag != "" && tag != "-" {
			name = tag
		} else {
			name = typeField.Name
		}

		// Add non-zero values to the update map
		if !field.IsZero() {
			updateData[name] = field.Interface()
		}
	}

	return updateData
}
