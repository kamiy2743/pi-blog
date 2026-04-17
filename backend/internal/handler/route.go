package handler

import (
	"net/http"

	"blog/internal/config"
	"blog/internal/di"
	createArticleHandler "blog/internal/handler/feature/admin/article/create"
	editArticleHandler "blog/internal/handler/feature/admin/article/edit"
	storeArticleHandler "blog/internal/handler/feature/admin/article/store"
	updateArticleHandler "blog/internal/handler/feature/admin/article/update"
	showAdminHandler "blog/internal/handler/feature/admin/show"
	"blog/internal/handler/inertia"
	"blog/internal/handler/middleware"

	"github.com/romsar/gonertia/v2"
)

func newMux(inertiaApp *gonertia.Inertia, container *di.Container) *http.ServeMux {
	mux := http.NewServeMux()

	setUpRoute(mux, inertiaApp, container)
	setUpArticleRoutes(mux, inertiaApp, container)
	setUpAdminRoutes(mux, inertiaApp, container)

	return mux
}

func setUpRoute(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
	container *di.Container,
) {
	mux.HandleFunc("GET /health", container.HealthHandler.Handle)

	mux.Handle("GET /", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			inertia.RenderNotFound(w, r, inertiaApp)
			return
		}
		container.ShowTopHandler.Handle(w, r)
	})))
}

func setUpArticleRoutes(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
	container *di.Container,
) {
	mux.Handle("GET /article", inertiaApp.Middleware(http.HandlerFunc(container.SearchArticleHandler.Handle)))

	mux.Handle("GET /article/{articleId}", inertiaApp.Middleware(http.HandlerFunc(container.ShowArticleHandler.Handle)))
}

func setUpAdminRoutes(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
	container *di.Container,
) {
	basicAuth := middleware.BasicAuth(
		"blog-admin",
		config.MustGetAdminBasicAuthUser(),
		config.MustGetAdminBasicAuthPass(),
	)
	handleAdmin := middleware.HandleWith(mux, basicAuth)

	handleAdmin("GET /admin", inertiaApp.Middleware(showAdminHandler.Handle(inertiaApp)))

	handleAdmin("GET /admin/article/new", inertiaApp.Middleware(createArticleHandler.Handle(inertiaApp)))
	handleAdmin("POST /admin/article/new", http.HandlerFunc(storeArticleHandler.Handle))

	handleAdmin("GET /admin/article/edit/{articleId}", inertiaApp.Middleware(http.HandlerFunc(editArticleHandler.Handle)))
	handleAdmin("POST /admin/article/edit/{articleId}", http.HandlerFunc(updateArticleHandler.Handle))
}
