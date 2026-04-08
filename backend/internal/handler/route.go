package handler

import (
	"net/http"

	"blog/internal/config"
	"blog/internal/di"
	"blog/internal/domain/article"
	createArticleHandler "blog/internal/handler/admin/article/create"
	editArticleHandler "blog/internal/handler/admin/article/edit"
	storeArticleHandler "blog/internal/handler/admin/article/store"
	updateArticleHandler "blog/internal/handler/admin/article/update"
	showAdminHandler "blog/internal/handler/admin/show"
	searchArticleHandler "blog/internal/handler/article/search"
	showArticleHandler "blog/internal/handler/article/show"
	healthHandler "blog/internal/handler/health"
	showNotFoundHandler "blog/internal/handler/notfound/show"
	"blog/internal/middleware"

	"github.com/romsar/gonertia/v2"
)

func newMux(inertiaApp *gonertia.Inertia, container *di.Container) *http.ServeMux {
	mux := http.NewServeMux()

	setUpRoute(mux, inertiaApp, container)
	setUpArticleRoutes(mux, inertiaApp)
	setUpAdminRoutes(mux, inertiaApp)

	return mux
}

func setUpRoute(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
	container *di.Container,
) {
	mux.HandleFunc("GET /health", healthHandler.Handle)

	mux.Handle("GET /", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			showNotFoundHandler.Handle(inertiaApp)(w, r)
			return
		}
		container.ShowTopHandler.Handle(w, r)
	})))
}

func setUpArticleRoutes(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
) {
	mux.Handle("GET /article", inertiaApp.Middleware(searchArticleHandler.Handle(inertiaApp)))

	mux.Handle("GET /article/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := article.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			showNotFoundHandler.Handle(inertiaApp)(w, r)
			return
		}
		showArticleHandler.Handle(inertiaApp, articleID)(w, r)
	})))
}

func setUpAdminRoutes(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
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

	handleAdmin("GET /admin/article/edit/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := article.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			showNotFoundHandler.Handle(inertiaApp)(w, r)
			return
		}
		editArticleHandler.Handle(inertiaApp, articleID)(w, r)
	})))
	handleAdmin("POST /admin/article/edit/{articleId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := article.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		updateArticleHandler.Handle(w, r, articleID)
	}))
}
