package main

import (
	"context"
	"log"
	"os"

	"blog/internal/db"
	"blog/internal/seed"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("使い方: blog-seed <seed_name>")
	}

	entClient, err := db.OpenEntClient()
	if err != nil {
		log.Fatalf("Ent client 初期化に失敗しました: %v", err)
	}
	defer entClient.Close()

	switch os.Args[1] {
	case "default":
		if err := seed.RunDefault(context.Background(), entClient); err != nil {
			log.Fatalf("default seed の投入に失敗しました: %v", err)
		}
	default:
		log.Fatalf("未対応の seed 名です: %s", os.Args[1])
	}
}
