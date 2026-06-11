package ws

import (
	"log"

	"gopkg.in/yaml.v3"
	"hei-gin/sdk/config"
)

// WSConfig holds WebSocket configuration, loaded from config.yaml's `ws:` section.
type WSConfig struct {
	ReadBufferSize          int `yaml:"read_buffer_size"`
	WriteBufferSize         int `yaml:"write_buffer_size"`
	HeartbeatInterval       int `yaml:"heartbeat_interval"`
	InstanceTTL             int `yaml:"instance_ttl"`
	StaleCleanInterval      int `yaml:"stale_clean_interval"`
	RateLimitWindow         int `yaml:"rate_limit_window"`
	RateLimitMax            int `yaml:"rate_limit_max"`
	DedupTTL                int `yaml:"dedup_ttl"`
	PollTimeout             int `yaml:"poll_timeout"`
	PongTimeout             int `yaml:"pong_timeout"`
	WriteTimeout            int `yaml:"write_timeout"`
	OnlineBroadcastInterval int `yaml:"online_broadcast_interval"`
}

// loadConfig reads ws config from config.C.Raw["ws"].
func loadConfig() WSConfig {
	cfg := defaultConfig()
	if config.C == nil {
		return cfg
	}
	raw, ok := config.C.Raw["ws"]
	if !ok {
		return cfg
	}
	data, err := yaml.Marshal(raw)
	if err != nil {
		log.Printf("[ws] failed to marshal ws config: %v", err)
		return cfg
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("[ws] failed to unmarshal ws config: %v", err)
		return cfg
	}
	return cfg
}

func defaultConfig() WSConfig {
	return WSConfig{
		ReadBufferSize:          1024,
		WriteBufferSize:         1024,
		HeartbeatInterval:       30,
		PongTimeout:             10,
		WriteTimeout:            10,
		OnlineBroadcastInterval: 60,
	}
}
