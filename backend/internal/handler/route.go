package handler

import (
	"net/http"

	"blog/internal/config"
	"blog/internal/di"
	createArticleHandler "blog/internal/handler/feature/admin/article/create"
	editArticleHandler "blog/internal/handler/feature/admin/article/edit"
	storeArticleHandler "blog/internal/handler/feature/admin/article/store"
	updateArticleHandler "blog/internal/handler/feature/admin/article/update"
	"blog/internal/handler/middleware"

	"github.com/romsar/gonertia/v2"
)

func newMux(inertiaApp *gonertia.Inertia, container *di.Container) *http.ServeMux {
	mux := http.NewServeMux()

	setUpRootRoutes(mux, inertiaApp, container)
	setUpArticleRoutes(mux, inertiaApp, container)
	setUpAdminRoutes(mux, inertiaApp, container)

	return mux
}

func setUpRootRoutes(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
	container *di.Container,
) {
	mux.Handle("GET /", InertiaPage(inertiaApp, container.ShowTopHandler.Handle))

	mux.HandleFunc("GET /health", container.HealthHandler.Handle)
}

func setUpArticleRoutes(
	mux *http.ServeMux,
	inertiaApp *gonertia.Inertia,
	container *di.Container,
) {
	mux.Handle("GET /article", inertiaApp.Middleware(http.HandlerFunc(container.SearchArticleHandler.Handle)))

	mux.Handle("GET /article/{articleId}", InertiaPage(inertiaApp, container.ShowArticleHandler.Handle))
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

	handleAdmin("GET /admin", inertiaApp.Middleware(http.HandlerFunc(container.ShowAdminHandler.Handle)))

	handleAdmin("GET /admin/article/new", inertiaApp.Middleware(createArticleHandler.Handle(inertiaApp)))
	handleAdmin("POST /admin/article/new", http.HandlerFunc(storeArticleHandler.Handle))

	handleAdmin("GET /admin/article/{articleId}", inertiaApp.Middleware(http.HandlerFunc(editArticleHandler.Handle)))
	handleAdmin("POST /admin/article/{articleId}", http.HandlerFunc(updateArticleHandler.Handle))

	handleAdmin("GET /admin/category", inertiaApp.Middleware(http.HandlerFunc(container.EditCategoryHandler.Handle)))
	handleAdmin("POST /admin/category", http.HandlerFunc(container.StoreCategoryHandler.Handle))
	handleAdmin("POST /admin/category/{categoryId}", http.HandlerFunc(container.UpdateCategoryHandler.Handle))
	handleAdmin("POST /admin/category/{categoryId}/delete", http.HandlerFunc(container.DestroyCategoryHandler.Handle))
}
