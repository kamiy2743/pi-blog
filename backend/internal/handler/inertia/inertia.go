package inertia

import (
	"errors"
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
		var displayableError *handlererror.DisplayableError
		if errors.As(err, &displayableError) {
			RenderError(w, r, inertiaApp, *displayableError)
			return
		}

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

type PrepareInputResult[T any] struct {
	Input              T
	HasValidationError bool
	ValidationErrors   []handlererror.ValidationError
	Request            *http.Request
	OK                 bool
}

func PrepareInput[T any](
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	toInput func(*http.Request) (T, []handlererror.ValidationError, *handlererror.DisplayableError),
) PrepareInputResult[T] {
	input, validationErrors, err := toInput(r)
	if err != nil {
		RenderError(w, r, inertiaApp, *err)
		return PrepareInputResult[T]{
			Input:   input,
			Request: r,
			OK:      false,
		}
	}

	return PrepareInputResult[T]{
		Input:              input,
		HasValidationError: len(validationErrors) > 0,
		ValidationErrors:   validationErrors,
		Request:            setValidationErrors(r, validationErrors),
		OK:                 true,
	}
}

func setValidationErrors(r *http.Request, validationErrors []handlererror.ValidationError) *http.Request {
	if len(validationErrors) == 0 {
		return r
	}

	errs := gonertia.ValidationErrors{}
	for _, validationError := range validationErrors {
		errs[validationError.Field] = validationError.Message
	}

	return r.WithContext(gonertia.SetValidationErrors(r.Context(), errs))
}

func triggerStaticErrorPage(w http.ResponseWriter) {
	http.Error(w, http.StatusText(staticErrorPageStatusCode), staticErrorPageStatusCode)
}
