package main

import (
	"log"
	"net/http"
	"os"

	"blog/internal/handler"
	"blog/internal/model"

	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal(".env に HOST が未設定です。")
	}
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
	buildDir := os.Getenv("BUILD_DIR")
	if buildDir == "" {
		log.Fatal(".env に BUILD_DIR が未設定です。")
	}

	inertiaOptions := []inertia.Option{
		inertia.WithSSR(frontURL),
	}
	inertiaApp, err := inertia.NewFromFile(rootTemplate, inertiaOptions...)
	if err != nil {
		log.Fatalf("Inertia のテンプレート読み込みに失敗しました: %v", err)
	}

	addr := host + ":" + port
	mux := http.NewServeMux()
	setupStaticRoutes(mux, buildDir)
	setupRootRoutes(mux, inertiaApp)
	setupArticleRoutes(mux, inertiaApp)
	setupAdminRoutes(mux, inertiaApp)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
}

func setupStaticRoutes(mux *http.ServeMux, buildDir string) {
	fileServer := http.FileServer(http.Dir(buildDir))
	mux.Handle("GET /build/", http.StripPrefix("/build/", fileServer))
	mux.Handle("HEAD /build/", http.StripPrefix("/build/", fileServer))
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
	mux.Handle("GET /admin", inertiaApp.Middleware(handler.ShowAdmin(inertiaApp)))
	mux.Handle("GET /admin/", inertiaApp.Middleware(handler.ShowAdmin(inertiaApp)))

	mux.Handle("GET /admin/article/new", inertiaApp.Middleware(handler.CreateArticle(inertiaApp)))
	mux.HandleFunc("POST /admin/article/new", handler.StoreArticle)

	mux.Handle("GET /admin/article/edit/{articleId}", inertiaApp.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		handler.EditArticle(inertiaApp, articleID)(w, r)
	})))
	mux.HandleFunc("POST /admin/article/edit/{articleId}", func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		handler.UpdateArticle(w, r, articleID)
	})

	mux.HandleFunc("POST /admin/article/publish/{articleId}", func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		handler.UpdatePublishSetting(w, r, articleID)
	})
}
