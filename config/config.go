package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	_ "time/tzdata"
)

type DatabaseConfig struct {
	DatabaseUri             string        `yaml:"database_uri" env:"DATABASE_URI"`
	DatabaseTimeout         time.Duration `yaml:"database_timeout" env:"DATABASE_TIMEOUT" env-default:"60s"`
	DatabaseMaxOpenConns    int           `env:"DATABASE_MAX_OPEN_CONNS" yaml:"database_max_open_conns" env-default:"5"`
	DatabaseMaxIdleConns    int           `env:"DATABASE_MAX_IDLE_CONNS" yaml:"database_max_idle_conns" env-default:"1"`
	DatabaseConnMaxLifetime time.Duration `env:"DATABASE_CONN_MAX_LIFETIME" yaml:"database_conn_max_lifetime" env-default:"15m"`
	DatabaseConnMaxIdletime time.Duration `env:"DATABASE_CONN_MAX_IDELTIME" yaml:"database_conn_max_ideltime" env-default:"5m"`
	DatabaseDebug           int           `yaml:"database_debug" env:"DATABASE_DEBUG"`
	DatabaseSlog            bool          `yaml:"database_slog" env:"DATABASE_SLOG"`
}

type CacheConfig struct {
<<<<<<< HEAD
=======
	// redis://<CACHE_USER>:<CACHE_PASSWORD>@<CACHE_END>:<CACHE_PORT>/<CACHE_DB>
>>>>>>> 9743c32 (chore: remove unused docker-entrypoint.sh)
	CacheUri            string        `yaml:"cache_uri" env:"CACHE_URI"`
	CacheUserSessionTTL int64         `env:"CACHE_USER_SESSION_TTL"`
	CachePoolSize       int           `env:"CACHE_POOL_SIZE"`
	CachePoolTimeout    time.Duration `env:"CACHE_POOL_TIMEOUT"`
}

<<<<<<< HEAD
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
=======
type StorageConfig struct {
	// s3://<STORAGE_ACCESS_KEY>:<STORAGE_SECRET_KEY>@<STORAGE_ENDPOINT>/<STORAGE_BUCKET_NAME>?sslmode=<STORAGE_USE_SSL:enable>
	StorageUri string `env:"STORAGE_URI"`
}

type Config struct {
	AppVersion      string         `env:"-"`
	JakartaLocation *time.Location `env:"-"`

	ServiceName      string   `env:"SERVICE_NAME" yaml:"service_name"`
	Env              string   `env:"ENV" yaml:"env"`
	JsonLog          bool     `yaml:"json_log" env:"JSON_LOG"`
	MsgLog           bool     `yaml:"msg_log" env:"MSG_LOG"`
	StartupMessage   bool     `env:"STARTUP_MESSAGE" yaml:"startup_message" env-default:"true"`
	MinioInitLog     bool     `env:"MINIO_INIT_LOG" yaml:"minio_init_log"`
	CorsAllowDomains []string `env:"CORS_ALLOW_DOMAINS"`
	FiberLog         bool     `env:"FIBER_LOG" yaml:"fiber_log"`
	SessionCookie    string   `env:"SESSION_COOKIE" env-default:"session_id"`
	ApiDocs          bool     `env:"API_DOCS" yaml:"api_docs" env-default:"true"`

	DB      DatabaseConfig
	Cache   CacheConfig
	Storage StorageConfig
>>>>>>> 9743c32 (chore: remove unused docker-entrypoint.sh)

	PublishUrl string `env:"PUBLISH_URL"`
	Host       string `yaml:"host" env:"HOST"`
	Port       int    `yaml:"port" env:"PORT"`

	DDService                string `env:"DD_SERVICE"`
	OTELExporterOTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	UptraceDSN               string `env:"UPTRACE_DSN"`
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
			loc, err := time.LoadLocation("Asia/Jakarta")
			if err == nil {
				conf.JakartaLocation = loc
			} else {
				slog.Error("Error loading Jakarta location", slog.Any("error", err))
			}

			envFile := ".env"
			slog.Info("ENV=" + os.Getenv("ENV"))
			if os.Getenv("ENV") == "test" {
				envFile = ".env.test.env"
			}

			if err := cleanenv.ReadConfig(envFile, conf); err != nil {
				if err := cleanenv.ReadEnv(conf); err != nil {
					slog.Error("Error reading env", slog.Any("error", err))
				} else {
					slog.Info("Reading config shell ENV")
				}
			} else {
				slog.Info("Reading config from " + envFile + " file")
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

func (c *Config) IsStaging() bool {
	return strings.HasPrefix(strings.ToLower(c.Env), "stag")
}

func (c *Config) IsDevelopment() bool {
	env := strings.ToLower(c.Env)
	return strings.HasPrefix(env, "dev") || env == "local"
}
