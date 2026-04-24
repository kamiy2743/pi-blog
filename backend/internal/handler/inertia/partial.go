package inertia

import (
	"net/http"
	"strings"
)

func ShouldIncludeProp(r *http.Request, component string, propName string) bool {
	partialComponent := r.Header.Get("X-Inertia-Partial-Component")
	partialData := r.Header.Get("X-Inertia-Partial-Data")

	if partialComponent == "" || partialData == "" {
		return true
	}
	if partialComponent != component {
		return false
	}

	for _, name := range strings.Split(partialData, ",") {
		if strings.TrimSpace(name) == propName {
			return true
		}
	}
	return false
}

func IsPartialReload(r *http.Request, component string) bool {
	return r.Header.Get("X-Inertia-Partial-Component") == component &&
		r.Header.Get("X-Inertia-Partial-Data") != ""
}
