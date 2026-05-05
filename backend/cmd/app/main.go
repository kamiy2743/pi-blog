package main

import (
	"log"
	"net/http"

	"blog/internal/config"
	"blog/internal/db"
	"blog/internal/handler"
	"blog/internal/handler/session"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	entClient, err := db.OpenEntClient()
	if err != nil {
		log.Fatalf("Ent client 初期化に失敗しました: %v", err)
	}
	defer entClient.Close()

	sessionManager := session.NewSessionManager(config.MustGetAppEnv())

	handler, err := handler.NewHTTPHandler(entClient, sessionManager)
	if err != nil {
		log.Fatalf("HTTP handler 初期化に失敗しました: %v", err)
	}

	addr := ":" + config.MustGetPort()
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
		return
	}
	log.Printf("listening on %s", addr)
}
