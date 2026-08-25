package config

import "errors"

type Config struct {
	GrpcAddr    string
	PostgresURL string
	KafkaBroker string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(httpAddr, postgresUrl, kafkaBroker string) (*Config, error) {
	if httpAddr == "" || postgresUrl == "" || kafkaBroker == "" {
		return nil, errLoadConfig
	}
	return &Config{GrpcAddr: httpAddr, PostgresURL: postgresUrl, KafkaBroker: kafkaBroker}, nil
}
