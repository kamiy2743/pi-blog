package db

import (
	"database/sql"
	"fmt"

	"blog/internal/config"
	"blog/internal/ent"

	_ "github.com/go-sql-driver/mysql"
)

func OpenMySQL() (*sql.DB, error) {
	return sql.Open("mysql", getDSN())
}

func OpenEntClient() (*ent.Client, error) {
	return ent.Open("mysql", getDSN())
}

func getDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.MustGetMySQLUser(),
		config.MustGetMySQLPassword(),
		config.MustGetMySQLHost(),
		config.MustGetMySQLPort(),
		config.MustGetMySQLDatabase(),
	)
}
