package helper

import (
	"net/http"
)

func genericError(err error) *Error {
	return &Error{Code: "generic", Message: err.Error()}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	e, ok := err.(*Error)
	if !ok {
		e = genericError(err)
	}

	WriteJSON(w, status, &ResponseBase{
		Status: ResponseStatusError,
		Error:  e,
	})
}
