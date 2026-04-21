package handler

import (
	"net/http"

	"blog/internal/handler/handlererror"
	"blog/internal/handler/session"

	"github.com/romsar/gonertia/v2"
)

type HandlerResult interface {
	isResult()
}

type PageResult struct {
	Component        string
	Props            gonertia.Props
	StatusCode       int
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}

func (PageResult) isResult() {}

type RedirectResult struct {
	To    string
	Flash *session.Flash
}

func (RedirectResult) isResult() {}

type RedirectBackResult struct {
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}

func (RedirectBackResult) isResult() {}

type PageOptions struct {
	StatusCode       int
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}

type RedirectOptions struct {
	Flash *session.Flash
}

type RedirectBackOptions struct {
	ValidationErrors []handlererror.ValidationError
	Flash            *session.Flash
}

func Page(component string, props gonertia.Props, options ...PageOptions) HandlerResult {
	option := PageOptions{
		StatusCode: http.StatusOK,
	}
	if len(options) > 0 {
		option = options[0]
	}

	return PageResult{
		Component:        component,
		Props:            props,
		StatusCode:       option.StatusCode,
		ValidationErrors: option.ValidationErrors,
		Flash:            option.Flash,
	}
}

func Redirect(to string, options ...RedirectOptions) HandlerResult {
	var option RedirectOptions
	if len(options) > 0 {
		option = options[0]
	}

	return RedirectResult{
		To:    to,
		Flash: option.Flash,
	}
}

func RedirectBack(options ...RedirectBackOptions) HandlerResult {
	var option RedirectBackOptions
	if len(options) > 0 {
		option = options[0]
	}

	return RedirectBackResult{
		ValidationErrors: option.ValidationErrors,
		Flash:            option.Flash,
	}
}
