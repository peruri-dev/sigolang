package testconfig

import "sigolang/config"

var testConf *config.Config = &config.Config{
	Env: "test",
	DB: config.DatabseConfig{
		DatabaseUri: "file::memory:?cache=shared",
	},
	Host: "sigolang-example.com",
	Port: 8888,
}

func ReloadTestConfig() *config.Config {
	return config.ReloadTestConfig(testConf)
}
