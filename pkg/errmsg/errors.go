package errmsg

import (
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func Errors[T any](err error, payloads ...*T) (code int, errors any) {
	var payload *T
	errors = make(map[string][]string)
	code = 500

	if len(payloads) > 0 {
		payload = payloads[0]
	}

	// REQUEST VALIDATION ERRORS
	if payload != nil {
		if errValidator, ok := err.(validator.ValidationErrors); ok {
			code, errors = errorValidationHandler(errValidator, payload)
		}
	}

	// DATABASE ERRORS
	if errPq, ok := err.(*pq.Error); ok {
		code, errors = errorPqHandler(errPq)
	}

	// CUSTOM ERRORS
	if errHttp, ok := err.(*CustomError); ok {
		code, errors = errorCustomHandler(errHttp)
	}

	return code, errors
}
