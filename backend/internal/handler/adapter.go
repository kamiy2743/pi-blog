package handler

import (
	"errors"
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
	handle func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError),
) http.Handler {
	return inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := handle(r)
		if err != nil {
			if pageResult, ok := result.(handlerresult.PageResult); ok {
				respondPartialPageError(w, r, inertiaApp, pageResult, err)
				return
			}
			respondPageError(w, r, inertiaApp, err)
			return
		}
		if result == nil {
			respondPageError(w, r, inertiaApp, errors.New("result が nil です"))
			return
		}

		switch typedResult := result.(type) {
		case handlerresult.PageResult:
			respondPageResult(w, r, inertiaApp, typedResult)
		case handlerresult.RedirectResult:
			respondRedirectResult(w, r, inertiaApp, typedResult)
		case handlerresult.RedirectBackResult:
			respondRedirectBackResult(w, r, inertiaApp, typedResult)
		default:
			respondPageError(w, r, inertiaApp, errors.New("未知の result 型です"))
		}
	}))
}

func InertiaAction(
	inertiaApp *gonertia.Inertia,
	handle func(*http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError),
) http.Handler {
	return inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := handle(r)
		if err != nil {
			respondActionError(w, r, inertiaApp, err)
			return
		}
		if result == nil {
			respondActionError(w, r, inertiaApp, errors.New("result が nil です"))
			return
		}

		switch typedResult := result.(type) {
		case handlerresult.RedirectResult:
			respondRedirectResult(w, r, inertiaApp, typedResult)
		case handlerresult.RedirectBackResult:
			respondRedirectBackResult(w, r, inertiaApp, typedResult)
		default:
			respondActionError(w, r, inertiaApp, errors.New("未知の result 型です"))
		}
	}))
}

func respondPageResult(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	result handlerresult.PageResult,
) {
	props := result.Props
	if props == nil {
		props = gonertia.Props{}
	}

	sessionPayload := popSessionPayload(r)
	validationErrors := selectValidationErrors(
		handlererror.ValidationErrorsToMap(result.ValidationErrors),
		sessionPayload.ValidationErrors,
	)
	flash := selectFlash(result.Flash, sessionPayload.Flash)

	props["validationErrors"] = validationErrors
	props["flash"] = session.FlashToMap(flash)

	inertia.Render(w, r, inertiaApp, result.StatusCode, result.Component, props)
}

func respondRedirectResult(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	result handlerresult.RedirectResult,
) {
	saveFlash(r, result.Flash)
	inertiaApp.Redirect(w, r, result.To, http.StatusSeeOther)
}

func respondRedirectBackResult(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	result handlerresult.RedirectBackResult,
) {
	saveValidationErrors(r, result.ValidationErrors)
	saveFlash(r, result.Flash)
	inertiaApp.Redirect(w, r, getRedirectBackURL(r), http.StatusSeeOther)
}

func respondPageError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	err error,
) {
	var displayableError *handlererror.DisplayableError
	if errors.As(err, &displayableError) {
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

func respondPartialPageError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	result handlerresult.PageResult,
	err error,
) {
	var displayableError *handlererror.DisplayableError
	if errors.As(err, &displayableError) && inertia.IsPartialReload(r, result.Component) {
		result.Flash = &session.Flash{Error: displayableError.Message}
		respondPageResult(w, r, inertiaApp, result)
		return
	}

	respondPageError(w, r, inertiaApp, err)
}

func respondActionError(
	w http.ResponseWriter,
	r *http.Request,
	inertiaApp *gonertia.Inertia,
	err error,
) {
	var displayableError *handlererror.DisplayableError
	if errors.As(err, &displayableError) {
		respondRedirectBackResult(w, r, inertiaApp, handlerresult.RedirectBackResult{
			Flash: &session.Flash{Error: displayableError.Message},
		})
		return
	}

	log.Print(err)
	respondRedirectBackResult(w, r, inertiaApp, handlerresult.RedirectBackResult{
		Flash: &session.Flash{Error: "エラーが発生しました。"},
	})
}

func saveValidationErrors(r *http.Request, validationErrors []handlererror.ValidationError) {
	manager, ok := session.SessionManagerFromContext(r.Context())
	if !ok {
		return
	}
	manager.SaveValidationErrors(r, validationErrors)
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

func selectValidationErrors(
	resultValidationErrors map[string]string,
	sessionValidationErrors map[string]string,
) map[string]string {
	if len(resultValidationErrors) > 0 {
		return resultValidationErrors
	}
	if len(sessionValidationErrors) > 0 {
		return sessionValidationErrors
	}
	return map[string]string{}
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
