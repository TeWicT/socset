package config

import "errors"

type Config struct {
	GRPCAddr                 string
	PostgresURL              string
	JWTSecret                string
	KafkaBroker              string
	OtelExporterOtlpEndpoint string
	OtelServiceName          string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(grpcAddr, postgresUrl, jwtSecret, kafkaBroker, otelExporterOtlpEndpoint, otelServiceName string) (*Config, error) {
	if grpcAddr == "" || postgresUrl == "" || len(jwtSecret) < 32 || kafkaBroker == "" || otelExporterOtlpEndpoint == "" || otelServiceName == "" {
		return nil, errLoadConfig
	}
	return &Config{GRPCAddr: grpcAddr, PostgresURL: postgresUrl, JWTSecret: jwtSecret, KafkaBroker: kafkaBroker, OtelExporterOtlpEndpoint: otelExporterOtlpEndpoint, OtelServiceName: otelServiceName}, nil
}
