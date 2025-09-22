package cmd

import (
	"fmt"
	"os"

	"sigolang/config"
	"sigolang/db/migrations"
	libcmd "sigolang/lib/cmd"
	"sigolang/lib/db"

	"github.com/spf13/cobra"
)

func AddDBCommands(rootCmd *cobra.Command) {
	getOpts := func() *libcmd.DBCommandsOptions {
		c := config.Get()

		dbConn, err := db.Open(&c.DB)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		return &libcmd.DBCommandsOptions{
			DBConn:     dbConn,
			Migrations: migrations.Migrations,
		}
	}

	dbCommands := libcmd.NewDBCommands(getOpts)
	rootCmd.AddCommand(
		dbCommands.Init,
		dbCommands.Migrate,
		dbCommands.Rollback,
		dbCommands.Lock,
		dbCommands.Unlock,
		dbCommands.CreateGo,
		dbCommands.CreateSQL,
		dbCommands.Status,
		dbCommands.MarkApplied,
	)
}
