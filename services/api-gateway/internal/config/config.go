package config

import "errors"

type Config struct {
	HttpAddr        string
	GRPCAddrAuth    string
	GRPCAddrProfile string
	RedisAddr       string
	JWTSecret       string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(httpAddr, grpcAddrAuth, grpcAddrProfile, redisAddr, jwtSecret string) (*Config, error) {
	if httpAddr == "" || grpcAddrAuth == "" || grpcAddrProfile == "" || redisAddr == "" || len(jwtSecret) < 32 {
		return nil, errLoadConfig
	}
	return &Config{HttpAddr: httpAddr, GRPCAddrAuth: grpcAddrAuth, GRPCAddrProfile: grpcAddrProfile, RedisAddr: redisAddr, JWTSecret: jwtSecret}, nil
}
