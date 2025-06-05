package cmd

import (
	"fmt"
	"os"

	"sigolang/config"
	"sigolang/db/seeds"
	"sigolang/lib/db"

	"github.com/spf13/cobra"
)

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed database in development",

	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		if !(c.IsDevelopment() || c.IsTesting()) {
			fmt.Println("seeding should be on `dev` or `test` env")
			os.Exit(1)
		}

		dbConn, err := db.Open(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		seeds.ResetSchema(cmd.Context(), dbConn)
	},
}
