package main

import (
	"log"
	"net/http"
	"os"

	"blog/internal/handler"
	"blog/internal/model"
)

func main() {
	mux := http.NewServeMux()
	setupRootRoutes(mux)
	setupArticleRoutes(mux)
	setupAdminRoutes(mux)

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal(".env に HOST が未設定です。")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal(".env に PORT が未設定です。")
	}
	addr := host + ":" + port

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func setupRootRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handler.Top)
	mux.HandleFunc("GET /health", handler.Health)
}

func setupArticleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /article/{articleId}", func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			handler.Top(w, r)
			return
		}
		handler.Article(w, r, articleID)
	})
}

func setupAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin", handler.ShowAdmin)
	mux.HandleFunc("GET /admin/article/new", handler.CreateArticle)
	mux.HandleFunc("POST /admin/article/new", handler.StoreArticle)
	mux.HandleFunc("GET /admin/article/edit/{articleId}", func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			handler.ShowAdmin(w, r)
			return
		}
		handler.EditArticle(w, r, articleID)
	})
	mux.HandleFunc("POST /admin/article/edit/{articleId}", func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			handler.ShowAdmin(w, r)
			return
		}
		handler.UpdateArticle(w, r, articleID)
	})
	mux.HandleFunc("POST /admin/article/publish/{articleId}", func(w http.ResponseWriter, r *http.Request) {
		articleID, err := model.ParseArticleID(r.PathValue("articleId"))
		if err != nil {
			handler.ShowAdmin(w, r)
			return
		}
		handler.UpdatePublishSetting(w, r, articleID)
	})
}
