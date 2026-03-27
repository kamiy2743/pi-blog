package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/internal/model"

	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	appEnvRaw := mustGetEnv("APP_ENV")
	appEnv, err := model.ParseAppEnv(appEnvRaw)
	if err != nil {
		log.Fatal(err)
	}
	port := mustGetEnv("PORT")
	frontURL := mustGetEnv("FRONT_URL")
	rootTemplate := mustGetEnv("INERTIA_ROOT_TEMPLATE")

	inertiaOptions := []inertia.Option{
		inertia.WithSSR(frontURL),
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
	faviconHref := mustGetEnv("TEMPLATE_FAVICON_HREF")
	cssHref := mustGetEnv("TEMPLATE_CSS_HREF")
	appScriptSrc := mustGetEnv("TEMPLATE_APP_SCRIPT_SRC")

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
	authUser := mustGetSecret("ADMIN_BASIC_AUTH_USER_FILE")
	authPass := mustGetSecret("ADMIN_BASIC_AUTH_PASS_FILE")

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

func mustGetEnv(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Fatalf(".env に %s が未設定です。", envName)
	}
	return value
}

func mustGetSecret(envName string) string {
	path := os.Getenv(envName)
	if path == "" {
		log.Fatalf(".env に %s が未設定です。", envName)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("%s の読み込みに失敗しました: %v", envName, err)
	}

	value := strings.TrimSpace(string(raw))
	if value == "" {
		log.Fatalf("%s の内容が空です。", envName)
	}
	return value
}
