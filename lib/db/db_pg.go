package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"sigolang/config"

	"github.com/uptrace/bun"

	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
	//sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

func init() {
	bunFactories = append(bunFactories, &BunFactory{
		Prefixes: []string{
			"postgres://", "postgresql://", "unix://",
		},
		Opener: func(c *config.DatabaseConfig) (db *bun.DB, err error) {
			dsn := c.DatabaseUri
			var dbConn *sql.DB
			connector := pgdriver.NewConnector(
				pgdriver.WithDSN(dsn),
				pgdriver.WithTimeout(time.Duration(c.DatabaseTimeout)*time.Second),
			)
			//if Datadog is configured, send sql traces there
			//if config.DatadogAgentUrl != "" {
			//	sqltrace.Register("postgres", pgdriver.Driver{}, sqltrace.WithServiceName("lndhub.go"))
			//	dbConn = sqltrace.OpenDB(connector)
			//} else {
			dbConn = sql.OpenDB(connector)
			//}
			db = bun.NewDB(dbConn, pgdialect.New(), bun.WithDiscardUnknownColumns())
			db.SetMaxOpenConns(c.DatabaseMaxOpenConns)
			db.SetMaxIdleConns(c.DatabaseMaxIdleConns)
			db.SetConnMaxLifetime(time.Duration(c.DatabaseConnMaxLifetime) * time.Minute)
			db.SetConnMaxIdleTime(time.Duration(c.DatabaseConnMaxIdletime) * time.Minute)
			db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("mydb")))

			ctx := context.Background()
			_, err = db.NewSelect().ColumnExpr("1").Exec(ctx)
			if err != nil {
				return nil, fmt.Errorf("error SELECT 1 postresql: %w", err)
			}

			slog.Info("postresql connected", slog.String("dsn", dsn))

			return db, nil
		},
	})
}
