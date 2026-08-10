package config

import "errors"

type Config struct {
	GRPCAddr    string
	PostgresURL string
	JWTSecret   string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(http_addr, postgres_url, jwtSecret string) (*Config, error) {
	if http_addr == "" || postgres_url == "" || len(jwtSecret) < 32 {
		return nil, errLoadConfig
	}
	return &Config{GRPCAddr: http_addr, PostgresURL: postgres_url, JWTSecret: jwtSecret}, nil
}
