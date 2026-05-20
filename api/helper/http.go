package helper

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Write(b)
}
