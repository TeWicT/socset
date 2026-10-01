package config

import "errors"

type Config struct {
	GRPCAddr    string
	PostgresURL string
	KafkaBroker string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(grpcAddr, postgresUrl, kafkaBroker string) (*Config, error) {
	if grpcAddr == "" || postgresUrl == "" || kafkaBroker == "" {
		return nil, errLoadConfig
	}
	return &Config{GRPCAddr: grpcAddr, PostgresURL: postgresUrl, KafkaBroker: kafkaBroker}, nil
}
