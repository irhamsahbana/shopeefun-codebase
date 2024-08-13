package errmsg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func errorValidationHandler[T any](err error, payload *T) (int, map[string][]string) {
	var (
		errorMessages = make(map[string][]string)
		code          = 400
	)

	for _, err := range err.(validator.ValidationErrors) {
		var (
			// Get the JSON tag name
			namespace  = err.Namespace()               // ex: UpdateInterestRequest.interest
			fieldParts = strings.Split(namespace, ".") // ex: [UpdateInterestRequest, interest]

			field      string
			fieldInMsg string
			message    string

			value     = err.Value()
			valueType = reflect.TypeOf(value)

			// Get the error message
		)
		lastField := fieldParts[len(fieldParts)-1]                                // get the last element
		fieldParts = fieldParts[1:]                                               // remove the first element
		field = strings.Join(fieldParts, ".")                                     // join the rest of the elements
		if strings.Contains(lastField, "_") && strings.Contains(lastField, "]") { // check if the last element contains "_" and "]", ex: interested_in[0]
			// fieldInMsg = field
			// remove characters between "[" and "]"
			fieldInMsg = strings.ReplaceAll(lastField, "_", " ")
			fieldInMsg = fieldInMsg[:strings.Index(fieldInMsg, "[")] // remove characters after "[" ("interested_in[0]" => "interested_in")
		} else {
			fieldInMsg = strings.ReplaceAll(lastField, "_", " ")
			if strings.Contains(fieldInMsg, "[") {
				fieldInMsg = fieldInMsg[:strings.Index(fieldInMsg, "[")] // remove characters after "[" ("interested_in[0]" => "interested_in")
			}
		}

		if err.Param() != "" {
			// message = fmt.Sprintf("field validation for '%s' failed on the '%s' tag with param '%s'", field, err.Tag(), err.Param())
			message = fmt.Sprintf("validasi untuk '%s' gagal pada tag '%s' dengan parameter '%s'", fieldInMsg, err.Tag(), err.Param())
		} else {
			// message = fmt.Sprintf("field validation for '%s' failed on the '%s' tag", field, err.Tag())
			message = fmt.Sprintf("validasi untuk '%s' gagal pada tag '%s'", fieldInMsg, err.Tag())
		}

		// get validate tag that causes the error
		switch err.Tag() {
		case "required":
			// message = fmt.Sprintf("%s is required.", fieldInMsg)
			message = fmt.Sprintf("%s harus diisi.", fieldInMsg)
		case "email":
			// message = fmt.Sprintf("%s is not a valid email address.", field)
			message = fmt.Sprintf("%s bukan alamat email yang valid.", fieldInMsg)
		case "email_blacklist":
			// message = fmt.Sprintf("email %v is not allowed.", value)
			message = fmt.Sprintf("email %v tidak diizinkan.", value)
		case "strong_password":
			// message = fmt.Sprintf("%s must be at least 12 characters and contain at least one uppercase letter, one lowercase letter, and one number.", fieldInMsg)
			message = fmt.Sprintf("%s minimal 12 karakter dan harus mengandung setidaknya satu huruf besar, satu huruf kecil, dan satu angka.", fieldInMsg)
		case "exist":
			// message = "resource is not exist."
			message = "sumber data tidak ditemukan."
		case "datetime":
			// message = fmt.Sprintf("%s is not a valid datetime format (Ex: %s).", fieldInMsg, err.Param())
			message = fmt.Sprintf("%s bukan format tanggal dan waktu yang valid (Contoh: %s).", fieldInMsg, err.Param())
		case "ulid":
			// message = fmt.Sprintf("%s is not a valid ULID.", fieldInMsg)
			message = fmt.Sprintf("%s bukan ULID yang valid.", fieldInMsg)
		case "base64":
			// message = fmt.Sprintf("%s is not a valid base64 format.", fieldInMsg)
			message = fmt.Sprintf("%s bukan format base64 yang valid.", fieldInMsg)
		case "base64url":
			// message = fmt.Sprintf("%s is not a valid base64url format.", fieldInMsg)
			message = fmt.Sprintf("%s bukan format base64url yang valid.", fieldInMsg)
		case "base64rawurl":
			// message = fmt.Sprintf("%s is not a valid base64rawurl format.", fieldInMsg)
			message = fmt.Sprintf("%s bukan format base64rawurl yang valid.", fieldInMsg)
		case "min":
			// check if the field is a number or a string
			if valueType.Kind() == reflect.Int || valueType.Kind() == reflect.Int8 || valueType.Kind() == reflect.Int16 || valueType.Kind() == reflect.Int32 || valueType.Kind() == reflect.Int64 || valueType.Kind() == reflect.Float32 || valueType.Kind() == reflect.Float64 {
				// message = fmt.Sprintf("%s must be at least %s.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus minimal %s.", fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.String {
				// message = fmt.Sprintf("%s must be at least %s characters.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus minimal %s karakter.", fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.Slice {
				// message = fmt.Sprintf("%s must have at least %s items.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus minimal %s item.", fieldInMsg, err.Param())
			}
		case "max":
			// check if the field is a number or a string
			if _, ok := value.(int); ok {
				// message = fmt.Sprintf("%s must not be greater than %s.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus tidak lebih dari %s.", fieldInMsg, err.Param())
			}
			if _, ok := value.(float64); ok {
				// message = fmt.Sprintf("%s must not be greater than %s.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus tidak lebih dari %s.", fieldInMsg, err.Param())
			}
			if _, ok := value.(string); ok {
				// message = fmt.Sprintf("%s must not be greater than %s characters.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus tidak lebih dari %s karakter.", fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.Slice {
				// message = fmt.Sprintf("%s must not have more than %s items.", fieldInMsg, err.Param())
				message = fmt.Sprintf("%s harus tidak lebih dari %s item.", fieldInMsg, err.Param())
			}
		case "gt":
			// message = fmt.Sprintf("%s must be greater than %s.", fieldInMsg, err.Param())
			message = fmt.Sprintf("%s harus lebih dari %s.", fieldInMsg, err.Param())
		case "gte":
			// message = fmt.Sprintf("%s must be greater than or equal to %s.", fieldInMsg, err.Param())
			message = fmt.Sprintf("%s harus lebih dari atau sama dengan %s.", fieldInMsg, err.Param())
		case "lt":
			// message = fmt.Sprintf("%s must be less than %s.", fieldInMsg, err.Param())
			message = fmt.Sprintf("%s harus kurang dari %s.", fieldInMsg, err.Param())
		case "lte":
			// message = fmt.Sprintf("%s must be less than or equal to %s.", fieldInMsg, err.Param())
			message = fmt.Sprintf("%s harus kurang dari atau sama dengan %s.", fieldInMsg, err.Param())
		case "latitude":
			// message = fmt.Sprintf("%s must be a valid latitude.", fieldInMsg)
			message = fmt.Sprintf("%s harus latitude yang valid.", fieldInMsg)
		case "longitude":
			// message = fmt.Sprintf("%s must be a valid longitude.", fieldInMsg)
			message = fmt.Sprintf("%s harus longitude yang valid.", fieldInMsg)
		case "numeric":
			// message = fmt.Sprintf("%s must be a number.", fieldInMsg)
			message = fmt.Sprintf("%s harus angka.", fieldInMsg)
		case "eqfield":
			eqField := err.Param()
			eqFieldName := ""
			eqFieldTag, _ := reflect.TypeOf(payload).Elem().FieldByName(eqField)
			eqFieldJSONTag := eqFieldTag.Tag.Get("json")
			eqFieldQueryTag := eqFieldTag.Tag.Get("query")
			eqFieldFormTag := eqFieldTag.Tag.Get("form")
			eqFieldParamsTag := eqFieldTag.Tag.Get("params")

			if eqFieldJSONTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldJSONTag, "_", " ")
			}
			if eqFieldQueryTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldQueryTag, "_", " ")
			}
			if eqFieldFormTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldFormTag, "_", " ")
			}
			if eqFieldParamsTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldParamsTag, "_", " ")
			}

			// message = fmt.Sprintf("%s must be equal to %s.", fieldInMsg, eqFieldName)
			message = fmt.Sprintf("%s harus sama dengan %s.", fieldInMsg, eqFieldName)
		case "oneof":
			// message = fmt.Sprintf("%s must be one of %s.", fieldInMsg, err.Param())
			// message = fmt.Sprintf("%s harus salah satu dari %s.", fieldInMsg, err.Param())

			// change param to be more readable
			// ex: "oneof=1 2 3" => "1, 2, atau 3"
			oneOfValues := strings.Split(err.Param(), " ")
			oneOfValues[len(oneOfValues)-1] = "atau " + oneOfValues[len(oneOfValues)-1]
			oneOfValuesStr := strings.Join(oneOfValues, ", ")
			message = fmt.Sprintf("%s harus salah satu dari %s.", fieldInMsg, oneOfValuesStr)
		case "unique_in_slice":
			// message = fmt.Sprintf("%s elements must be unique.", fieldInMsg)
			message = fmt.Sprintf("elemen %s harus unik.", fieldInMsg)
		}

		errorMessages[field] = append(errorMessages[field], message)
	}

	return code, errorMessages
}
