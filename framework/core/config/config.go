package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 是进程启动配置根（YAML / 环境变量）。
//
// Author: Charlie
type Config struct {
	App     AppConfig     `mapstructure:"app"`
	DB      DBConfig      `mapstructure:"db"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Auth    AuthConfig    `mapstructure:"auth"`
	CORS    CORSConfig    `mapstructure:"cors"`
	XxlJob  XxlJobConfig  `mapstructure:"xxl_job"`
	Storage StorageConfig `mapstructure:"storage"`
	IDGen   IDGenConfig   `mapstructure:"id_generator"`
	Modules ModulesConfig `mapstructure:"modules"`
	Audit   AuditConfig   `mapstructure:"audit"`
}

// AppConfig 应用基础信息与监听地址。
//
// Author: Charlie
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
	Timezone string `mapstructure:"timezone"`
}

// DBConfig 数据库连接与连接池。
//
// Author: Charlie
type DBConfig struct {
	URL         string `mapstructure:"url"`
	Echo        bool   `mapstructure:"echo"`
	PoolSize    int    `mapstructure:"pool_size"`
	MaxOverflow int    `mapstructure:"max_overflow"`
}

// RedisConfig Redis 连接。
//
// Author: Charlie
type RedisConfig struct {
	URL            string `mapstructure:"url"`
	MaxConnections int    `mapstructure:"max_connections"`
}

// AuthConfig 会话、登录保护与 Cookie 策略。
//
// Author: Charlie
type AuthConfig struct {
	TokenName                   string   `mapstructure:"token_name"`
	TokenTTLSeconds             int      `mapstructure:"token_ttl_seconds"`
	TokenTTLShortSeconds        int      `mapstructure:"token_ttl_short_seconds"`
	PortalRegisterEnabled       bool     `mapstructure:"portal_register_enabled"`
	LoginFailureWindowSeconds   int      `mapstructure:"login_failure_window_seconds"`
	LoginAccountMaxFailures     int      `mapstructure:"login_account_max_failures"`
	LoginIPMaxFailures          int      `mapstructure:"login_ip_max_failures"`
	LoginLockSeconds            int      `mapstructure:"login_lock_seconds"`
	CaptchaTTLSeconds           int      `mapstructure:"captcha_ttl_seconds"`
	PasswordCryptoKeyTTLSeconds int      `mapstructure:"password_crypto_key_ttl_seconds"`
	SessionIdleTimeoutSeconds   int      `mapstructure:"session_idle_timeout_seconds"`
	SessionBindIP               bool     `mapstructure:"session_bind_ip"`
	SessionBindUserAgent        bool     `mapstructure:"session_bind_user_agent"`
	MaxConcurrentSessions       int      `mapstructure:"max_concurrent_sessions"`
	SessionCookieEnabled        bool     `mapstructure:"session_cookie_enabled"`
	SessionCookieName           string   `mapstructure:"session_cookie_name"`
	SessionCookieSecure         bool     `mapstructure:"session_cookie_secure"`
	SessionCookieSameSite       string   `mapstructure:"session_cookie_samesite"`
	SessionCookiePath           string   `mapstructure:"session_cookie_path"`
	DefaultPassword             string   `mapstructure:"default_password"`
	AuthWhitelist               []string `mapstructure:"auth_whitelist"`
}

// CORSConfig 跨域策略。
//
// Author: Charlie
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
}

// XxlJobConfig XXL-JOB 执行器（嵌在 API 进程）。
//
// Author: Charlie
type XxlJobConfig struct {
	Enabled     bool               `mapstructure:"enabled"`
	AccessToken string             `mapstructure:"access_token"`
	Admin       XxlJobAdminConfig  `mapstructure:"admin"`
	Executor    XxlJobExecutorConf `mapstructure:"executor"`
}

// XxlJobAdminConfig 调度中心地址。
//
// Author: Charlie
type XxlJobAdminConfig struct {
	Addresses string `mapstructure:"addresses"`
}

// XxlJobExecutorConf 执行器注册名与端口。
//
// Author: Charlie
type XxlJobExecutorConf struct {
	AppName string `mapstructure:"appname"`
	Port    int    `mapstructure:"port"`
}

// StorageConfig 对象存储（local / S3 兼容）参数。
//
// Author: Charlie
type StorageConfig struct {
	Provider             string `mapstructure:"provider"`
	Bucket               string `mapstructure:"bucket"`
	Endpoint             string `mapstructure:"endpoint"`
	AccessKey            string `mapstructure:"access_key"`
	SecretKey            string `mapstructure:"secret_key"`
	Region               string `mapstructure:"region"`
	UseSSL               bool   `mapstructure:"use_ssl"`
	PresignExpireSeconds int    `mapstructure:"presign_expire_seconds"`
	BaseURL              string `mapstructure:"base_url"`
	PublicPath           string `mapstructure:"public_path"`
	LocalRoot            string `mapstructure:"local_root"`
}

// IDGenConfig 雪花 ID 的 worker / datacenter。
//
// Author: Charlie
type IDGenConfig struct {
	WorkerID     int64 `mapstructure:"worker_id"`
	DatacenterID int64 `mapstructure:"datacenter_id"`
}

// ModulesConfig 模块启用 / 禁用过滤。
//
// Author: Charlie
type ModulesConfig struct {
	Disabled []string `mapstructure:"disabled"`
	Enabled  []string `mapstructure:"enabled"` // 非空时仅运行这些模块
}

// AuditConfig 操作审计队列容量。
//
// Author: Charlie
type AuditConfig struct {
	OperationQueueSize int `mapstructure:"operation_queue_size"`
}

// Load 从指定 YAML 路径加载配置，并叠加 HEI_ 环境变量。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("HEI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "hei-gin")
	v.SetDefault("app.host", "127.0.0.1")
	v.SetDefault("app.port", 8000)
	v.SetDefault("app.debug", true)
	v.SetDefault("app.timezone", "Asia/Shanghai")
	v.SetDefault("db.url", "postgres://postgres:123456@127.0.0.1:5432/hei_gin?sslmode=disable")
	v.SetDefault("db.pool_size", 10)
	v.SetDefault("db.max_overflow", 20)
	v.SetDefault("redis.url", "redis://:123456@127.0.0.1:6379/4")
	v.SetDefault("redis.max_connections", 100)
	v.SetDefault("auth.token_name", "Authorization")
	v.SetDefault("auth.token_ttl_seconds", 14400)
	v.SetDefault("auth.token_ttl_short_seconds", 7200)
	v.SetDefault("auth.portal_register_enabled", true)
	v.SetDefault("auth.login_failure_window_seconds", 900)
	v.SetDefault("auth.login_account_max_failures", 5)
	v.SetDefault("auth.login_ip_max_failures", 30)
	v.SetDefault("auth.login_lock_seconds", 900)
	v.SetDefault("auth.captcha_ttl_seconds", 300)
	v.SetDefault("auth.password_crypto_key_ttl_seconds", 600)
	v.SetDefault("auth.session_cookie_enabled", true)
	v.SetDefault("auth.session_cookie_name", "Authorization")
	v.SetDefault("auth.session_cookie_samesite", "lax")
	v.SetDefault("auth.session_cookie_path", "/")
	v.SetDefault("auth.max_concurrent_sessions", 5)
	v.SetDefault("auth.session_bind_ip", true)
	v.SetDefault("auth.default_password", "123456")
	v.SetDefault("cors.allow_origins", []string{
		"http://localhost:5173", "http://127.0.0.1:5173",
		"http://localhost:5174", "http://127.0.0.1:5174",
		"http://localhost:5163", "http://127.0.0.1:5163",
	})
	v.SetDefault("cors.allow_credentials", true)
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allow_headers", []string{"Authorization", "Content-Type", "X-Request-Id", "Accept", "Origin", "X-Requested-With", "X-HEI-CSRF"})
	v.SetDefault("xxl_job.enabled", true)
	v.SetDefault("xxl_job.access_token", "default_token")
	v.SetDefault("xxl_job.admin.addresses", "http://127.0.0.1:9004/xxl-job-admin")
	v.SetDefault("xxl_job.executor.appname", "hei-gin-api")
	v.SetDefault("xxl_job.executor.port", 9999)
	v.SetDefault("storage.provider", "local")
	v.SetDefault("storage.bucket", "hei-gin")
	v.SetDefault("storage.public_path", "/api/v1/files")
	v.SetDefault("storage.local_root", "storage")
	v.SetDefault("storage.presign_expire_seconds", 3600)
	v.SetDefault("id_generator.worker_id", 1)
	v.SetDefault("id_generator.datacenter_id", 1)
	v.SetDefault("audit.operation_queue_size", 1000)
}

// Addr 返回 HTTP 监听地址 host:port。
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}
