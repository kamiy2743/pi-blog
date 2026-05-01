package handlerrequest

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONForm[T any](r *http.Request, dst *T) {
	_ = json.NewDecoder(r.Body).Decode(dst)
}
