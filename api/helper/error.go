package helper

import (
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, &ResponseBase{
		Status: ResponseStatusError,
		Error:  err,
	})
}
