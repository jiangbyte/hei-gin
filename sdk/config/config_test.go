package config

import "testing"

func TestValidateRuntimeRequiresCriticalFields(t *testing.T) {
	cfg := &Config{}
	if err := cfg.ValidateRuntime(true); err == nil {
		t.Fatal("ValidateRuntime should reject empty config")
	}
}

func TestValidateRuntimeAcceptsCompleteConfig(t *testing.T) {
	cfg := &Config{
		App:   AppConfig{Host: "127.0.0.1", Port: 8080},
		DB:    DatabaseConfig{Host: "127.0.0.1", Port: 3306, User: "root", Database: "hei"},
		Redis: RedisConfig{Host: "127.0.0.1", Port: 6379},
		Token: TokenConfig{ExpireSeconds: 3600, TokenName: "Authorization"},
	}
	if err := cfg.ValidateRuntime(true); err != nil {
		t.Fatalf("ValidateRuntime returned error: %v", err)
	}
}

func TestValidateMigrationRequiresDBFields(t *testing.T) {
	cfg := &Config{}
	if err := cfg.ValidateMigration(); err == nil {
		t.Fatal("ValidateMigration should reject empty config")
	}
}
