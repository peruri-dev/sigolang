package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabseConfig struct {
	DatabaseUri     string `yaml:"database_uri" env:"DATABASE_URI"`
	DatabaseTimeout int    `yaml:"database_timeout" env:"DATABASE_TIMEOUT"`
}

type Config struct {
	Env      string `env:"ENV" yaml:"env"`
	CacheUri string `yaml:"cache_uri" env:"CACHE_URI"`

	DB DatabseConfig

	Host string `yaml:"host" env:"HOST"`
	Port int    `yaml:"port" env:"PORT"`
}

var conf *Config
var once sync.Once

func Get() *Config {
	if conf != nil {
		return conf
	}

	once.Do(
		func() {
			conf = &Config{}

			if err := cleanenv.ReadConfig(".env", conf); err != nil {
				fmt.Printf("Error reading config file, %s\n", err)
				if err := cleanenv.ReadEnv(conf); err != nil {
					fmt.Printf("Error reading env, %s\n", err)
				}
			}
		})

	return conf
}

func ReloadTestConfig(testConf *Config) *Config {
	conf = testConf
	return conf
}

func (c *Config) IsProduction() bool {
	return strings.HasPrefix(strings.ToLower(c.Env), "prod")
}

func (c *Config) IsTesting() bool {
	return strings.HasPrefix(strings.ToLower(c.Env), "test")
}

func (c *Config) IsDevelopment() bool {
	return strings.HasPrefix(strings.ToLower(c.Env), "dev")
}
