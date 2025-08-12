package config

import (
	"encoding/json"
	"os"
)

// TunnelMapping defines the relationship between a local listener and a remote target.
type TunnelMapping struct {
	LocalAddress  string `json:"local_address"`
	RemoteAddress string `json:"remote_address"`
}

// ClientConfig holds the configuration for the client application.
type ClientConfig struct {
	ServerAddress string          `json:"server_address"`
	Password      string          `json:"password"`
	Tunnels       []TunnelMapping `json:"tunnels"`
}

// ServerConfig holds the configuration for the server application.
type ServerConfig struct {
	ListenAddress string `json:"listen_address"`
	Password      string `json:"password"`
}

// LoadClientConfig reads and parses the client configuration file.
func LoadClientConfig(path string) (*ClientConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ClientConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadServerConfig reads and parses the server configuration file.
func LoadServerConfig(path string) (*ServerConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
