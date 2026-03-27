package main

import (
	"log"
	"net/http"

	"blog/internal/config"
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/internal/model"

	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	appEnv := config.MustGetAppEnv()
	port := config.MustGetPort()
	ssrURL := config.MustGetSSRURL()
	rootTemplate := config.MustGetInertiaRootTemplate()

	inertiaOptions := []inertia.Option{
		inertia.WithSSR(ssrURL),
	}
	inertiaApp, err := inertia.NewFromFile(rootTemplate, inertiaOptions...)
	if err != nil {
		log.Fatalf("Inertia のテンプレート読み込みに失敗しました: %v", err)
	}
	configureTemplateAssets(appEnv, inertiaApp)

	addr := ":" + port
	mux := http.NewServeMux()
	handler := middleware.Chain(
		http.NewCrossOriginProtection().Handler(mux),
		middleware.NormalizePath(),
	)

	setupRootRoutes(mux, inertiaApp)
	setupArticleRoutes(mux, inertiaApp)
	setupAdminRoutes(mux, inertiaApp)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
		return
	}
	log.Printf("listening on %s", addr)
}

func configureTemplateAssets(appEnv model.AppEnv, inertiaApp *inertia.Inertia) {
	faviconHref := config.MustGetTemplateFaviconHref()
	cssHref := config.MustGetTemplateCSSHref()
	appScriptSrc := config.MustGetTemplateAppScriptSrc()

	inertiaApp.ShareTemplateData("faviconHref", faviconHref)
	inertiaApp.ShareTemplateData("cssHref", cssHref)
	inertiaApp.ShareTemplateData("useViteClient", appEnv == model.AppEnvDev)
	inertiaApp.ShareTemplateData("appScriptSrc", appScriptSrc)
}

func setupRootRoutes(mux *http.ServeMux, inertiaApp *inertia.Inertia) {
	mux.HandleFunc("GET /health", handler.Health)

	mux.Handle("GET /", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handler.ShowNotFound(inertiaApp)(w, r)
			return
		}
		handler.ShowTop(inertiaApp)(w, r)
	})))
}

func setupArticleRoutes(mux *http.ServeMux, inertiaApp *inertia.Inertia) {
	mux.Handle("GET /article", inertiaApp.Middleware(handler.ShowArticleList(inertiaApp)))

	mux.Handle("GET /article/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			handler.ShowNotFound(inertiaApp)(w, r)
			return
		}
		handler.ShowArticle(inertiaApp, articleID)(w, r)
	})))
}

func setupAdminRoutes(mux *http.ServeMux, inertiaApp *inertia.Inertia) {
	authUser := config.MustGetAdminBasicAuthUser()
	authPass := config.MustGetAdminBasicAuthPass()

	basicAuth := middleware.BasicAuth("blog-admin", authUser, authPass)
	handleAdmin := middleware.HandleWith(mux, basicAuth)

	handleAdmin("GET /admin", inertiaApp.Middleware(handler.ShowAdmin(inertiaApp)))

	handleAdmin("GET /admin/article/new", inertiaApp.Middleware(handler.CreateArticle(inertiaApp)))
	handleAdmin("POST /admin/article/new", http.HandlerFunc(handler.StoreArticle))

	handleAdmin("GET /admin/article/edit/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			handler.ShowNotFound(inertiaApp)(w, r)
			return
		}
		handler.EditArticle(inertiaApp, articleID)(w, r)
	})))
	handleAdmin("POST /admin/article/edit/{articleId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		handler.UpdateArticle(w, r, articleID)
	}))

	handleAdmin("POST /admin/article/publish/{articleId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		handler.UpdatePublishSetting(w, r, articleID)
	}))
}
