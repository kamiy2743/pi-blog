package main

import (
	"log"
	"net/http"
	"os"

	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/internal/model"

	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal(".env に PORT が未設定です。")
	}
	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		log.Fatal(".env に FRONT_URL が未設定です。")
	}
	rootTemplate := os.Getenv("INERTIA_ROOT_TEMPLATE")
	if rootTemplate == "" {
		log.Fatal(".env に INERTIA_ROOT_TEMPLATE が未設定です。")
	}

	inertiaOptions := []inertia.Option{
		inertia.WithSSR(frontURL),
	}
	inertiaApp, err := inertia.NewFromFile(rootTemplate, inertiaOptions...)
	if err != nil {
		log.Fatalf("Inertia のテンプレート読み込みに失敗しました: %v", err)
	}

	addr := ":" + port
	mux := http.NewServeMux()
	handler := http.NewCrossOriginProtection().Handler(mux)
	setupRootRoutes(mux, inertiaApp)
	setupArticleRoutes(mux, inertiaApp)
	setupAdminRoutes(mux, inertiaApp)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
		return
	}
	log.Printf("listening on %s", addr)
}

func setupRootRoutes(mux *http.ServeMux, inertiaApp *inertia.Inertia) {
	mux.HandleFunc("GET /health", handler.Health)

	mux.Handle("GET /", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handler.ShowTop(inertiaApp)(w, r)
	})))
}

func setupArticleRoutes(mux *http.ServeMux, inertiaApp *inertia.Inertia) {
	mux.Handle("GET /article", inertiaApp.Middleware(handler.ShowArticleList(inertiaApp)))
	mux.Handle("GET /article/", inertiaApp.Middleware(handler.ShowArticleList(inertiaApp)))

	mux.Handle("GET /article/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		handler.ShowArticle(inertiaApp, articleID)(w, r)
	})))
}

func setupAdminRoutes(mux *http.ServeMux, inertiaApp *inertia.Inertia) {
	authUser := os.Getenv("ADMIN_BASIC_AUTH_USER")
	if authUser == "" {
		log.Fatal(".env に ADMIN_BASIC_AUTH_USER が未設定です。")
	}
	authPass := os.Getenv("ADMIN_BASIC_AUTH_PASS")
	if authPass == "" {
		log.Fatal(".env に ADMIN_BASIC_AUTH_PASS が未設定です。")
	}

	basicAuth := middleware.BasicAuth("blog-admin", authUser, authPass)
	handleAdmin := middleware.HandleWith(mux, basicAuth)

	handleAdmin("GET /admin", inertiaApp.Middleware(handler.ShowAdmin(inertiaApp)))
	handleAdmin("GET /admin/", inertiaApp.Middleware(handler.ShowAdmin(inertiaApp)))

	handleAdmin("GET /admin/article/new", inertiaApp.Middleware(handler.CreateArticle(inertiaApp)))
	handleAdmin("POST /admin/article/new", http.HandlerFunc(handler.StoreArticle))

	handleAdmin("GET /admin/article/edit/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
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
