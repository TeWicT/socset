package config

import "errors"

type Config struct {
	HttpAddr     string
	GRPCAddrAuth string
	JWTSecret    string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(httpAddr, grpcAddrAuth, jwtSecret string) (*Config, error) {
	if httpAddr == "" || grpcAddrAuth == "" || len(jwtSecret) < 32 {
		return nil, errLoadConfig
	}
	return &Config{HttpAddr: httpAddr, GRPCAddrAuth: grpcAddrAuth, JWTSecret: jwtSecret}, nil
}
