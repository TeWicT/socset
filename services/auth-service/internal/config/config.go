package config

import "errors"

type Config struct {
	GRPCAddr    string
	PostgresURL string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(http_addr, postgres_url string) (*Config, error) {
	if http_addr == "" || postgres_url == "" {
		return nil, errLoadConfig
	}
	return &Config{GRPCAddr: http_addr, PostgresURL: postgres_url}, nil
}
