package config

import "errors"

type Config struct {
	GRPCAddr                 string
	PostgresURL              string
	KafkaBroker              string
	OtelExporterOtlpEndpoint string
	OtelServiceName          string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(grpcAddr, postgresUrl, kafkaBroker, otelExporterOtlpEndpoint, otelServiceName string) (*Config, error) {
	if grpcAddr == "" || postgresUrl == "" || kafkaBroker == "" || otelExporterOtlpEndpoint == "" || otelServiceName == "" {
		return nil, errLoadConfig
	}
	return &Config{GRPCAddr: grpcAddr, PostgresURL: postgresUrl, KafkaBroker: kafkaBroker, OtelExporterOtlpEndpoint: otelExporterOtlpEndpoint, OtelServiceName: otelServiceName}, nil
}
