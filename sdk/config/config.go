package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App       AppConfig       `yaml:"app"`
	DB        DatabaseConfig  `yaml:"db"`
	Redis     RedisConfig     `yaml:"redis"`
	Token     TokenConfig     `yaml:"token"`
	SM2       SM2Config       `yaml:"sm2"`
	Auth      AuthConfig      `yaml:"auth"`
	CORS      CORSConfig      `yaml:"cors"`
	User      UserConfig      `yaml:"user"`
	Snowflake SnowflakeConfig `yaml:"snowflake"`
	Swagger   SwaggerConfig   `yaml:"swagger"`
	Raw       map[string]any  `yaml:",inline"`
}

type AppConfig struct {
	Name             string `yaml:"name"`
	Version          string `yaml:"version"`
	Env              string `yaml:"env"`
	Debug            bool   `yaml:"debug"`
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	UploadMaxSize    int64  `yaml:"upload_max_size"`
	TimeoutKeepAlive int    `yaml:"timeout_keep_alive"`
}

type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Database       string `yaml:"database"`
	PoolSize       int    `yaml:"pool_size"`
	MaxOverflow    int    `yaml:"max_overflow"`
	PoolRecycle    int    `yaml:"pool_recycle"`
	PoolPrePing    bool   `yaml:"pool_pre_ping"`
	PoolTimeout    int    `yaml:"pool_timeout"`
	ConnectTimeout int    `yaml:"connect_timeout"`
	Echo           bool   `yaml:"echo"`
}

type RedisConfig struct {
	Host                 string `yaml:"host"`
	Port                 int    `yaml:"port"`
	Password             string `yaml:"password"`
	Database             int    `yaml:"database"`
	MaxConnections       int    `yaml:"max_connections"`
	SocketConnectTimeout int    `yaml:"socket_connect_timeout"`
	SocketTimeout        int    `yaml:"socket_timeout"`
	RetryOnTimeout       bool   `yaml:"retry_on_timeout"`
	HealthCheckInterval  int    `yaml:"health_check_interval"`
}

type TokenConfig struct {
	ExpireSeconds int    `yaml:"expire_seconds"`
	TokenName     string `yaml:"token_name"`
}

type SM2Config struct {
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
}

type UserConfig struct {
	ResetPassword string `yaml:"reset_password"`
}

// SwaggerConfig holds Swagger API documentation settings.
type SwaggerConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AuthConfig struct {
	// PublicPaths lists routes that bypass authentication entirely.
	// Each entry is matched as a prefix against the full request path.
	PublicPaths             []string `yaml:"public_paths"`
	BusinessRegisterEnabled bool     `yaml:"business_register_enabled"`
}

type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

type SnowflakeConfig struct {
	Instance int64 `yaml:"instance"`
}

var ErrConfigNotFound = errors.New("config.yaml not found")

var C *Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	C = &Config{}
	return yaml.Unmarshal(data, C)
}

// FindAndLoad searches for config.yaml in the CWD, then walks up the directory tree.
func FindAndLoad() error {
	paths := []string{
		os.Getenv("HEI_CONFIG"),
		"config.yaml",
		"../config.yaml",
		"../../config.yaml",
	}
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			log.Printf("[Config] Loading from %s", p)
			return Load(p)
		}
	}
	// Find repo root by looking for go.work
	wd, _ := os.Getwd()
	dir := wd
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, "config.yaml")
		if _, err := os.Stat(candidate); err == nil {
			log.Printf("[Config] Loading from %s", candidate)
			return Load(candidate)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ErrConfigNotFound
}

func (c *Config) ValidateRuntime(redisRequired bool) error {
	if c == nil {
		return errors.New("config is nil")
	}

	var missing []string
	requireString(&missing, "app.host", c.App.Host)
	requirePositive(&missing, "app.port", c.App.Port)
	requireString(&missing, "db.host", c.DB.Host)
	requirePositive(&missing, "db.port", c.DB.Port)
	requireString(&missing, "db.user", c.DB.User)
	requireString(&missing, "db.database", c.DB.Database)
	requirePositive(&missing, "token.expire_seconds", c.Token.ExpireSeconds)
	requireString(&missing, "token.token_name", c.Token.TokenName)

	if redisRequired {
		requireString(&missing, "redis.host", c.Redis.Host)
		requirePositive(&missing, "redis.port", c.Redis.Port)
	}

	if len(missing) > 0 {
		return fmt.Errorf("invalid runtime config, missing or invalid fields: %s", strings.Join(missing, ", "))
	}
	return nil
}

func (c *Config) ValidateMigration() error {
	if c == nil {
		return errors.New("config is nil")
	}

	var missing []string
	requireString(&missing, "db.host", c.DB.Host)
	requirePositive(&missing, "db.port", c.DB.Port)
	requireString(&missing, "db.user", c.DB.User)
	requireString(&missing, "db.database", c.DB.Database)
	if len(missing) > 0 {
		return fmt.Errorf("invalid migration config, missing or invalid fields: %s", strings.Join(missing, ", "))
	}
	return nil
}

func requireString(missing *[]string, key, value string) {
	if strings.TrimSpace(value) == "" {
		*missing = append(*missing, key)
	}
}

func requirePositive(missing *[]string, key string, value int) {
	if value <= 0 {
		*missing = append(*missing, key)
	}
}
