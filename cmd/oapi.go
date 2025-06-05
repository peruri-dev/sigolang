package cmd

import (
	"fmt"

	"sigolang/config"
	"sigolang/internal/handler"
	"sigolang/internal/service"
	"sigolang/lib/transport"

	"github.com/spf13/cobra"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Print the OpenAPI spec",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()
		f := transport.InitFiber(c)
		svc := &service.Services{}
		api := handler.RegisterRoutes(f, svc)

		b, err := api.OpenAPI().YAML()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))

	},
}
