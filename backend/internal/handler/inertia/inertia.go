package inertia

import (
	"log"
	"net/http"
	"net/http/httptest"

	"blog/internal/handler/handlererror"

	"github.com/romsar/gonertia/v2"
)

const staticErrorPageStatusCode = http.StatusTeapot

func Render(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	component string,
	props gonertia.Props,
) {
	if err := inertiaApp.Render(w, r, component, props); err != nil {
		log.Print(err)
		triggerStaticErrorPage(w)
	}
}

func RenderWithStatus(
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
		triggerStaticErrorPage(w)
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

func RenderError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	displayableError handlererror.DisplayableError,
) {
	log.Print(displayableError.Err)

	RenderWithStatus(w, r, inertiaApp, displayableError.StatusCode, "ErrorPage", gonertia.Props{
		"statusCode":  displayableError.StatusCode,
		"statusText":  http.StatusText(displayableError.StatusCode),
		"message":     displayableError.Message,
		"description": displayableError.Description,
	})
}

func triggerStaticErrorPage(w http.ResponseWriter) {
	http.Error(w, http.StatusText(staticErrorPageStatusCode), staticErrorPageStatusCode)
}
