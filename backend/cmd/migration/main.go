package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blog/internal/db"
	"blog/internal/db/ent/migrate"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("使い方: blog-migrate {up|down|refresh}")
	}

	entClient, err := db.OpenEntClient()
	if err != nil {
		log.Fatalf("Ent client 初期化に失敗しました: %v", err)
	}
	defer entClient.Close()

	switch os.Args[1] {
	case "up":
		if err := entClient.Schema.Create(context.Background()); err != nil {
			log.Fatalf("Ent schema migration に失敗しました: %v", err)
		}
	case "down":
		if err := dropAllTables(); err != nil {
			log.Fatalf("Ent schema drop に失敗しました: %v", err)
		}
	case "refresh":
		if err := dropAllTables(); err != nil {
			log.Fatalf("Ent schema drop に失敗しました: %v", err)
		}
		if err := entClient.Schema.Create(context.Background()); err != nil {
			log.Fatalf("Ent schema migration に失敗しました: %v", err)
		}
	default:
		log.Fatalf("未対応の migration コマンドです: %s", os.Args[1])
	}
}

func dropAllTables() error {
	conn, err := db.OpenMySQL()
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}
	defer func() {
		_, _ = conn.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}()

	for i := len(migrate.Tables) - 1; i >= 0; i-- {
		query := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", migrate.Tables[i].Name)
		if _, err := conn.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
