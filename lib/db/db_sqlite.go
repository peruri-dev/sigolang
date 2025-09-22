package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"sigolang/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func init() {
	bunFactories = append(bunFactories, &BunFactory{
		Prefixes: []string{
			"file:",
		},
		Opener: func(c *config.DatabaseConfig) (db *bun.DB, err error) {
			dsn := c.DatabaseUri

			dbConn, err := sql.Open(sqliteshim.ShimName, dsn) // "file::memory:?cache=shared"
			if err != nil {
				return nil, err
			}
			db = bun.NewDB(dbConn, sqlitedialect.New(), bun.WithDiscardUnknownColumns())

			ctx := context.Background()
			_, err = db.NewSelect().ColumnExpr("1").Exec(ctx)
			if err != nil {
				return nil, fmt.Errorf("error SELECT 1 sqlite: %w", err)
			}

			slog.Info("sqlite connected", slog.String("dsn", dsn))

			return db, nil
		},
	})
}
