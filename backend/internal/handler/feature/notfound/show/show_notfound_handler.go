package show

import (
	"net/http"
	"net/http/httptest"

	"github.com/romsar/gonertia/v2"
)

func Handle(i *gonertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := renderWithStatus(w, http.StatusNotFound, func(target http.ResponseWriter) error {
			return i.Render(target, r, "NotFound", nil)
		}); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
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
