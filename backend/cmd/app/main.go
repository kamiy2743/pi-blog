package main

import (
	"log"
	"net/http"

	"blog/internal/config"
	"blog/internal/db"
	"blog/internal/di"
	"blog/internal/domain"
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

func main() {
	appEnv := config.MustGetAppEnv()
	port := config.MustGetPort()
	ssrURL := config.MustGetSSRURL()
	rootTemplate := config.MustGetInertiaRootTemplate()

	inertiaOptions := []gonertia.Option{
		gonertia.WithSSR(ssrURL),
	}
	inertiaApp, err := gonertia.NewFromFile(rootTemplate, inertiaOptions...)
	if err != nil {
		log.Fatalf("Inertia のテンプレート読み込みに失敗しました: %v", err)
	}
	configureTemplateAssets(appEnv, inertiaApp)

	entClient, err := db.OpenEntClient()
	if err != nil {
		log.Fatalf("Ent client 初期化に失敗しました: %v", err)
	}
	defer entClient.Close()

	container := di.NewContainer(entClient, inertiaApp)

	addr := ":" + port
	mux := http.NewServeMux()
	handler := middleware.Chain(
		http.NewCrossOriginProtection().Handler(mux),
		middleware.NormalizePath(),
	)

	setupRootRoutes(mux, inertiaApp, container)
	setupArticleRoutes(mux, inertiaApp)
	setupAdminRoutes(mux, inertiaApp)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
		return
	}
	log.Printf("listening on %s", addr)
}

func configureTemplateAssets(appEnv domain.AppEnv, inertiaApp *gonertia.Inertia) {
	faviconHref := config.MustGetTemplateFaviconHref()
	cssHref := config.MustGetTemplateCSSHref()
	appScriptSrc := config.MustGetTemplateAppScriptSrc()

	inertiaApp.ShareTemplateData("faviconHref", faviconHref)
	inertiaApp.ShareTemplateData("cssHref", cssHref)
	inertiaApp.ShareTemplateData("useViteClient", appEnv == domain.AppEnvDev)
	inertiaApp.ShareTemplateData("appScriptSrc", appScriptSrc)
}

func setupRootRoutes(
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
		container.ShowTopHandler().Handle(w, r)
	})))
}

func setupArticleRoutes(mux *http.ServeMux, inertiaApp *gonertia.Inertia) {
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

func setupAdminRoutes(mux *http.ServeMux, inertiaApp *gonertia.Inertia) {
	authUser := config.MustGetAdminBasicAuthUser()
	authPass := config.MustGetAdminBasicAuthPass()

	basicAuth := middleware.BasicAuth("blog-admin", authUser, authPass)
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
