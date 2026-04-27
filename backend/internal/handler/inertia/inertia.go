package inertia

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/romsar/gonertia/v3"
)

const staticErrorPageStatusCode = http.StatusTeapot

func Render(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	statusCode int,
	component string,
	props gonertia.Props,
) {
	recorder := httptest.NewRecorder()
	if err := inertiaApp.Render(recorder, r, component, props); err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(staticErrorPageStatusCode), staticErrorPageStatusCode)
		return
	}

	for key, values := range recorder.Header() {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(recorder.Body.Bytes())
}
