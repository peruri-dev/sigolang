package config

import (
	"log/slog"
	"sigolang/lib/util"
	"strings"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	DatabaseUri             string `yaml:"database_uri" env:"DATABASE_URI"`
	DatabaseTimeout         int    `yaml:"database_timeout" env:"DATABASE_TIMEOUT"`
	DatabaseMaxOpenConns    int    `env:"DATABASE_MAX_OPEN_CONNS" yaml:"database_max_open_conns" env-default:"5"`
	DatabaseMaxIdleConns    int    `env:"DATABASE_MAX_IDLE_CONNS" yaml:"database_max_idle_conns" env-default:"1"`
	DatabaseConnMaxLifetime int    `env:"DATABASE_CONN_MAX_LIFETIME" yaml:"database_conn_max_lifetime" env-default:"15"`
	DatabaseConnMaxIdletime int    `env:"DATABASE_CONN_MAX_IDELTIME" yaml:"database_conn_max_ideltime" env-default:"5"`
	DatabaseDebug           int    `yaml:"database_debug" env:"DATABASE_DEBUG"`
	DatabaseSlog            bool   `yaml:"database_slog" env:"DATABASE_SLOG"`
}

type CacheConfig struct {
	CacheUri            string        `yaml:"cache_uri" env:"CACHE_URI"`
	CacheUserSessionTTL int64         `env:"CACHE_USER_SESSION_TTL"`
	CachePoolSize       int           `env:"CACHE_POOL_SIZE"`
	CachePoolTimeout    time.Duration `env:"CACHE_POOL_TIMEOUT"`
}

type Config struct {
	AppVersion     string `env:"-"`
	ServiceName    string `env:"SERVICE_NAME" yaml:"service_name"`
	Env            string `env:"ENV" yaml:"env"`
	JsonLog        bool   `yaml:"json_log" env:"JSON_LOG"`
	ViteDist       bool   `env:"VITE_DIST" env-default:"false"`
	VitePort       int    `env:"VITE_PORT" env-default:"5173"`
	MsgLog         bool   `yaml:"msg_log" env:"MSG_LOG"`
	StartupMessage bool   `env:"STARTUP_MESSAGE" yaml:"startup_message"`

	DB    DatabaseConfig
	Cache CacheConfig

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

			if conf.ServiceName == "" {
				conf.ServiceName = util.GetExecutableName()
			}
		},
	)

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
