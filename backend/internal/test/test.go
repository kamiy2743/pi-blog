package test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"blog/internal/db"
	"blog/internal/di"
	"blog/internal/ent"
	"blog/internal/ent/enttest"
	"blog/internal/ent/migrate"
	"blog/internal/handler"

	_ "github.com/go-sql-driver/mysql"
)

type InitResult struct {
	EntClient *ent.Client
	Server    *httptest.Server
}

func Init(t *testing.T, containerOptions ...*di.ContainerOptions) InitResult {
	t.Helper()

	if err := resetTestDatabase(t.Context()); err != nil {
		t.Fatalf("mysql-test の初期化に失敗しました: %v", err)
	}

	entClient := enttest.Open(t, "mysql", db.GetDSN())
	t.Cleanup(func() {
		_ = entClient.Close()
	})

	httpHandler, err := handler.NewHTTPHandler(entClient, containerOptions...)
	if err != nil {
		t.Fatalf("HTTP ハンドラーの初期化に失敗しました: %v", err)
	}
	server := httptest.NewServer(httpHandler)
	t.Cleanup(func() {
		server.Close()
	})

	return InitResult{
		EntClient: entClient,
		Server:    server,
	}
}

func resetTestDatabase(ctx context.Context) error {
	conn, err := db.OpenMySQL()
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}
	defer func() {
		_, _ = conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
	}()

	for i := len(migrate.Tables) - 1; i >= 0; i-- {
		query := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", migrate.Tables[i].Name)
		if _, err := conn.ExecContext(ctx, query); err != nil {
			return err
		}
	}

	return nil
}
