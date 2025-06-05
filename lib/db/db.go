package db

import (
	"fmt"
	"strings"

	"sigolang/config"

	"github.com/uptrace/bun"

	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	//sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

type BunFactory struct {
	Prefixes []string
	Opener   func(c *config.Config) (*bun.DB, error)
}

var bunFactories []*BunFactory = []*BunFactory{}

func allPrefixes() string {
	prefixes := []string{}
	for _, bunFactory := range bunFactories {
		prefixes = append(prefixes, bunFactory.Prefixes...)
	}
	return strings.Join(prefixes, "|")
}

func Open(c *config.Config) (db *bun.DB, err error) {
	dsn := c.DB.DatabaseUri

	if dsn == "" {
		fmt.Println("not using database")
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
		db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("sigolang-db")))

		db.AddQueryHook(bundebug.NewQueryHook(
			// disable the hook
			bundebug.WithEnabled(false),
			// BUNDEBUG=1 logs failed queries
			// BUNDEBUG=2 logs all queries
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	return db, nil
}
