package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func HandleWith(mux *http.ServeMux, middlewares ...Middleware) func(string, http.Handler) {
	return func(pattern string, handler http.Handler) {
		mux.Handle(pattern, Chain(handler, middlewares...))
	}
}
