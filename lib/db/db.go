package db

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"sigolang/config"

	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/attribute"

	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/bun/extra/bunslog"
)

type BunFactory struct {
	Prefixes []string
	Opener   func(c *config.DatabaseConfig) (*bun.DB, error)
}

var bunFactories []*BunFactory = []*BunFactory{}

func allPrefixes() string {
	prefixes := []string{}
	for _, bunFactory := range bunFactories {
		prefixes = append(prefixes, bunFactory.Prefixes...)
	}
	return strings.Join(prefixes, "|")
}

func Open(c *config.DatabaseConfig) (db *bun.DB, err error) {
	dsn := c.DatabaseUri

	if dsn == "" {
		return
	}

	found := false

	for _, bunFactory := range bunFactories {
		for _, prefix := range bunFactory.Prefixes {
			if found = strings.HasPrefix(dsn, prefix); found {
				break
			}
		}

		if found {
			db, err = bunFactory.Opener(c)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("invalid database connection string %s, only (%s) is supported", dsn, allPrefixes())
	}

	if db != nil {
		dbname := "sigolang_dev"
		attrs := []attribute.KeyValue{}

		u, err := url.Parse(dsn)
		if err == nil {
			if len(u.Path) > 1 {
				dbname = u.Path[1:]
			}
			attrs = append(attrs, attribute.String("db.host", u.Host))
		}

		db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(dbname),
			bunotel.WithAttributes(attrs...)))

		if c.DatabaseSlog {
			db.AddQueryHook(bunslog.NewQueryHook(
				bunslog.WithQueryLogLevel(slog.LevelDebug),
				bunslog.WithSlowQueryLogLevel(slog.LevelWarn),
				bunslog.WithErrorQueryLogLevel(slog.LevelError),
				bunslog.WithSlowQueryThreshold(3*time.Second),
			))
		} else {
			db.AddQueryHook(bundebug.NewQueryHook(
				bundebug.WithEnabled(c.DatabaseDebug >= 1),
				bundebug.WithVerbose(c.DatabaseDebug >= 2),
			))
		}
	}

	return db, nil
}
