package testconfig

import (
	"log/slog"
	"path/filepath"
	"sigolang/config"
	"sigolang/lib/util"

	"github.com/ilyakaznacheev/cleanenv"
)

var testConf *config.Config = &config.Config{
	Env: "test",
	DB: config.DatabaseConfig{
		DatabaseUri: "file::memory:?cache=shared",
	},
	Host: "sigolang-example.com",
	Port: 8888,
}

func ReloadTestConfig() *config.Config {
	proj, _ := util.GetCurrentProjectRoot()
	path := filepath.Join(proj, ".env.test.env")

	if err := cleanenv.ReadConfig(path, testConf); err != nil {
		slog.Warn("fail read test environment", slog.String("path", path))
	}

	return config.ReloadTestConfig(testConf)
}
