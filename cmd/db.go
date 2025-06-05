package cmd

import (
	"fmt"
	"os"
	"strings"

	"sigolang/config"
	"sigolang/db/migrations"
	"sigolang/lib/db"

	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

var dbInitCmd = &cobra.Command{
	Use:   "db:init",
	Short: "create migration tables",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		err = migrator.Init(cmd.Context())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

var dbMigrateCmd = &cobra.Command{
	Use:   "db:migrate",
	Short: "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		group, err := migrator.Migrate(cmd.Context())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		if group.ID == 0 {
			fmt.Println("there are no new migrations to run")
			return
		}

		fmt.Printf("migrated to %s\n", group)
	},
}

var dbRollbackCmd = &cobra.Command{
	Use:   "db:rollback",
	Short: "rollback the last migration group",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		group, err := migrator.Rollback(cmd.Context())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		if group.ID == 0 {
			fmt.Println("there are no groups to roll back")
			return
		}

		fmt.Printf("rolled back %s\n", group)
	},
}

var dbLockCmd = &cobra.Command{
	Use:   "db:lock",
	Short: "lock migrations",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		err = migrator.Lock(cmd.Context())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
	},
}

var dbUnlockCmd = &cobra.Command{
	Use:   "db:unlock",
	Short: "unlock migrations",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		err = migrator.Unlock(cmd.Context())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
	},
}

var dbCreateGoCmd = &cobra.Command{
	Use:   "db:create_go",
	Short: "create Go migration",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		name := strings.Join(args, "_")
		mf, err := migrator.CreateGoMigration(cmd.Context(), name)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
		fmt.Printf("created migration %s (%s)", mf.Name, mf.Path)
	},
}

var dbCreateSqlCmd = &cobra.Command{
	Use:   "db:create_sql",
	Short: "create up and down SQL migrations",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		name := strings.Join(args, "_")
		files, err := migrator.CreateSQLMigrations(cmd.Context(), name)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		for _, mf := range files {
			fmt.Printf("created migration %s (%s)", mf.Name, mf.Path)
		}
	},
}

var dbStatusCmd = &cobra.Command{
	Use:   "db:status",
	Short: "print migrations status",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		ms, err := migrator.MigrationsWithStatus(cmd.Context())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		fmt.Printf("migrations: %s\n", ms)
		fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
		fmt.Printf("last migration group: %s\n", ms.LastGroup())
	},
}

var dbMarkAppliedCmd = &cobra.Command{
	Use:   "db:mark_applied",
	Short: "mark migrations as applied without actually running them",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		migrator := migrate.NewMigrator(dbConn, migrations.Migrations)

		group, err := migrator.Migrate(cmd.Context(), migrate.WithNopMigration())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		if group.ID == 0 {
			fmt.Println("there are no new migrations to mark as applied")
			return
		}

		fmt.Printf("marked as applied %s\n", group)
	},
}
