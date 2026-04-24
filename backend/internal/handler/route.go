package handler

import (
	"net/http"

	"blog/internal/config"
	"blog/internal/di"
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
	mux.Handle("GET /article", InertiaPage(inertiaApp, container.SearchArticleHandler.Handle))

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

	handleAdmin("GET /admin", InertiaPage(inertiaApp, container.ShowAdminHandler.Handle))

	handleAdmin("GET /admin/article/new", InertiaPage(inertiaApp, container.CreateArticleHandler.Handle))
	handleAdmin("POST /admin/article/new", InertiaAction(inertiaApp, container.StoreArticleHandler.Handle))

	handleAdmin("GET /admin/article/{articleId}", InertiaPage(inertiaApp, container.EditArticleHandler.Handle))
	handleAdmin("POST /admin/article/{articleId}", InertiaAction(inertiaApp, container.UpdateArticleHandler.Handle))

	handleAdmin("GET /admin/category", InertiaPage(inertiaApp, container.EditCategoryHandler.Handle))
	handleAdmin("POST /admin/category", InertiaAction(inertiaApp, container.StoreCategoryHandler.Handle))
	handleAdmin("POST /admin/category/{categoryId}", InertiaAction(inertiaApp, container.UpdateCategoryHandler.Handle))
	handleAdmin("POST /admin/category/{categoryId}/delete", InertiaAction(inertiaApp, container.DestroyCategoryHandler.Handle))
}
