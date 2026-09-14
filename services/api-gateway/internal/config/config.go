package config

import "errors"

type Config struct {
	HttpAddr        string
	GRPCAddrAuth    string
	GRPCAddrProfile string
	JWTSecret       string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(httpAddr, grpcAddrAuth, grpcAddrProfile, jwtSecret string) (*Config, error) {
	if httpAddr == "" || grpcAddrAuth == "" || grpcAddrProfile == "" || len(jwtSecret) < 32 {
		return nil, errLoadConfig
	}
	return &Config{HttpAddr: httpAddr, GRPCAddrAuth: grpcAddrAuth, GRPCAddrProfile: grpcAddrProfile, JWTSecret: jwtSecret}, nil
}
