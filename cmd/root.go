package cmd

import (
	"sigolang/config"
	"sigolang/internal/app"
	"sigolang/lib/util"

	"github.com/danielgtaylor/huma/v2/humacli"
)

type Options struct {
	Debug bool   `doc:"Enable debug logging"`
	Host  string `doc:"Hostname to listen on."`
	Port  int    `doc:"Port to listen on." short:"p"`
}

func applyOptions(opts *Options) *config.Config {
	c := config.Get()

	if opts.Host != "" {
		c.Host = opts.Host
	} else if c.Host == "" {
		c.Host = "localhost"
	}

	if opts.Port != 0 {
		c.Port = opts.Port
	} else if c.Port == 0 {
		c.Port = 3000
	}

	return c
}

func Execute(appVersion string) {
	// Then, create the CLI.
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		c := applyOptions(opts)
		a := app.NewApp(c)
		hooks.OnStart(a.Start)
		hooks.OnStop(a.Stop)
	})

	rootCmd := cli.Root()
	rootCmd.Use = util.GetExecutablePath()
	rootCmd.Version = appVersion

	AddDBCommands(rootCmd)
	rootCmd.AddCommand(dbSeedCmd)
	rootCmd.AddCommand(openapiCmd)

	// Run the thing!
	cli.Run()
}
