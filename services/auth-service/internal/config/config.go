package config

import "errors"

type Config struct {
	GRPCAddr    string
	PostgresURL string
	JWTSecret   string
	KafkaBroker string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(grpcAddr, postgresUrl, jwtSecret, kafkaBroker string) (*Config, error) {
	if grpcAddr == "" || postgresUrl == "" || len(jwtSecret) < 32 || kafkaBroker == "" {
		return nil, errLoadConfig
	}
	return &Config{GRPCAddr: grpcAddr, PostgresURL: postgresUrl, JWTSecret: jwtSecret, KafkaBroker: kafkaBroker}, nil
}
