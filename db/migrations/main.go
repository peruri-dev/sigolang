package migrations

import (
	"embed"
	"fmt"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := Migrations.DiscoverCaller(); err != nil {
		fmt.Println(err.Error())
	}

	if err := Migrations.Discover(sqlMigrations); err != nil {
		fmt.Println(err.Error())
	}
}
