package handler

import (
	"log"
	"net/http"

	"blog/internal/handler/handlererror"
	"blog/internal/handler/handlerresult"
	"blog/internal/handler/inertia"
	"blog/internal/handler/session"

	"github.com/romsar/gonertia/v3"
)

func InertiaPage(
	inertiaApp *gonertia.Inertia,
	handle func(*http.Request) (handlerresult.PageResult, error),
) http.Handler {
	return inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := handle(r)
		if err != nil {
			if validationError, ok := handlererror.AsValidationError(err); ok {
				respondPageResult(w, r, inertiaApp, result, validationError, nil)
				return
			}
			if inertia.IsPartialReload(r, result.Component) {
				if displayableError, ok := handlererror.AsDisplayableError(err); ok {
					respondPageResult(w, r, inertiaApp, result, &handlererror.ValidationError{}, &session.Flash{
						Error: displayableError.Message,
					})
					return
				}
			}
			respondPageError(w, r, inertiaApp, err)
			return
		}
		respondPageResult(w, r, inertiaApp, result, nil, nil)
	}))
}

func InertiaAction(
	inertiaApp *gonertia.Inertia,
	handle func(*http.Request) (handlerresult.ActionResult, error),
) http.Handler {
	return inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := handle(r)
		if err != nil {
			respondActionError(w, r, inertiaApp, err)
			return
		}
		respondActionResult(w, r, inertiaApp, result)
	}))
}

func respondPageResult(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	result handlerresult.PageResult,
	resultValidationError *handlererror.ValidationError,
	resultFlash *session.Flash,
) {
	props := result.Props
	if props == nil {
		props = gonertia.Props{}
	}

	sessionPayload := popSessionPayload(r)
	validationError := selectValidationError(
		resultValidationError,
		sessionPayload.ValidationError,
	)
	flash := selectFlash(resultFlash, sessionPayload.Flash)

	if validationError != nil && !validationError.IsEmpty() {
		props["validationErrors"] = validationError.Messages
	}
	if flash != nil && !flash.IsEmpty() {
		props["flash"] = session.FlashToMap(flash)
	}

	inertia.Render(w, r, inertiaApp, http.StatusOK, result.Component, props)
}

func respondActionResult(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	result handlerresult.ActionResult,
) {
	inertiaApp.Redirect(w, r, result.RedirectTo, http.StatusSeeOther)
}

func respondPageError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	err error,
) {
	if displayableError, ok := handlererror.AsDisplayableError(err); ok {
		inertia.RenderError(w, r, inertiaApp, *displayableError)
		return
	}

	log.Print(err)
	inertia.Render(w, r, inertiaApp, http.StatusInternalServerError, "ErrorPage", gonertia.Props{
		"statusCode":  http.StatusInternalServerError,
		"statusText":  http.StatusText(http.StatusInternalServerError),
		"message":     "エラーが発生しました。",
		"description": "時間をおいてから、もう一度お試しください。",
	})
}

func respondActionError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	err error,
) {
	if validationError, ok := handlererror.AsValidationError(err); ok {
		saveValidationError(r, validationError)
		respondRedirectBack(w, r, inertiaApp, nil)
		return
	}

	if displayableError, ok := handlererror.AsDisplayableError(err); ok {
		respondRedirectBack(w, r, inertiaApp, &session.Flash{Error: displayableError.Message})
		return
	}

	log.Print(err)
	respondRedirectBack(w, r, inertiaApp, &session.Flash{Error: "エラーが発生しました。"})
}

func respondRedirectBack(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	flash *session.Flash,
) {
	saveFlash(r, flash)
	inertiaApp.Redirect(w, r, getRedirectBackURL(r), http.StatusSeeOther)
}

func saveValidationError(r *http.Request, validationError *handlererror.ValidationError) {
	manager, ok := session.SessionManagerFromContext(r.Context())
	if !ok {
		return
	}
	manager.SaveValidationError(r, validationError)
}

func saveFlash(r *http.Request, flash *session.Flash) {
	manager, ok := session.SessionManagerFromContext(r.Context())
	if !ok {
		return
	}
	manager.SaveFlash(r, flash)
}

func popSessionPayload(r *http.Request) session.SessionPayload {
	manager, ok := session.SessionManagerFromContext(r.Context())
	if !ok {
		return session.SessionPayload{}
	}
	return manager.PopSessionPayload(r)
}

func getRedirectBackURL(r *http.Request) string {
	referer := r.Referer()
	if referer != "" {
		return referer
	}
	return "/"
}

func selectValidationError(
	resultValidationError *handlererror.ValidationError,
	sessionValidationError *handlererror.ValidationError,
) *handlererror.ValidationError {
	if resultValidationError != nil && !resultValidationError.IsEmpty() {
		return resultValidationError
	}
	if sessionValidationError != nil && !sessionValidationError.IsEmpty() {
		return sessionValidationError
	}
	return &handlererror.ValidationError{Messages: map[string]string{}}
}

func selectFlash(
	resultFlash *session.Flash,
	sessionFlash *session.Flash,
) *session.Flash {
	if resultFlash != nil && !resultFlash.IsEmpty() {
		return resultFlash
	}
	if sessionFlash != nil && !sessionFlash.IsEmpty() {
		return sessionFlash
	}
	return nil
}
