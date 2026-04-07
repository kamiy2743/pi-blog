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

func RenderError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	displayableError handlererror.DisplayableError,
) {
	log.Print(displayableError.Err)

	if err := renderWithStatus(w, displayableError.StatusCode, func(target http.ResponseWriter) error {
		return inertiaApp.Render(target, r, "ErrorPage", gonertia.Props{
			"statusCode":  displayableError.StatusCode,
			"statusText":  http.StatusText(displayableError.StatusCode),
			"message":     displayableError.Message,
			"description": displayableError.Description,
		})
	}); err != nil {
		log.Print(err)
		triggerStaticErrorPage(w)
	}
}

func renderWithStatus(
	w http.ResponseWriter,
	statusCode int,
	render func(target http.ResponseWriter) error,
) error {
	recorder := httptest.NewRecorder()
	if err := render(recorder); err != nil {
		return err
	}

	for key, values := range recorder.Header() {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(recorder.Body.Bytes())
	return nil
}

func triggerStaticErrorPage(w http.ResponseWriter) {
	http.Error(w, http.StatusText(staticErrorPageStatusCode), staticErrorPageStatusCode)
}
