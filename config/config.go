package config

import (
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabseConfig struct {
	DatabaseUri     string `yaml:"database_uri" env:"DATABASE_URI"`
	DatabaseTimeout int    `yaml:"database_timeout" env:"DATABASE_TIMEOUT"`
}

type Cache struct {
	CacheUri            string        `yaml:"cache_uri" env:"CACHE_URI"`
	CacheUserSessionTTL int64         `env:"CACHE_USER_SESSION_TTL"`
	CachePoolSize       int           `env:"CACHE_POOL_SIZE"`
	CachePoolTimeout    time.Duration `env:"CACHE_POOL_TIMEOUT"`
}

type Config struct {
	Env     string `env:"ENV" yaml:"env"`
	JsonLog bool   `yaml:"json_log" env:"JSON_LOG"`
	MsgLog  bool   `yaml:"msg_log" env:"MSG_LOG"`

	DB    DatabseConfig
	Cache Cache

	PublishUrl string `env:"PUBLISH_URL"`
	Host       string `yaml:"host" env:"HOST"`
	Port       int    `yaml:"port" env:"PORT"`
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
				if err := cleanenv.ReadEnv(conf); err != nil {
					slog.Error("Error reading env", slog.Any("error", err))
				} else {
					slog.Info("Reading config ENV")
				}
			} else {
				slog.Info("Reading config from .env file")
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
