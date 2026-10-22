package config

import "errors"

type Config struct {
	GrpcAddr                 string
	PostgresURL              string
	KafkaBroker              string
	OtelExporterOtlpEndpoint string
	OtelServiceName          string
}

var errLoadConfig = errors.New("err load config")

func CreateConfig(httpAddr, postgresUrl, kafkaBroker, otelExporterOtlpEndpoint, otelServiceName string) (*Config, error) {
	if httpAddr == "" || postgresUrl == "" || kafkaBroker == "" || otelExporterOtlpEndpoint == "" || otelServiceName == "" {
		return nil, errLoadConfig
	}
	return &Config{GrpcAddr: httpAddr, PostgresURL: postgresUrl, KafkaBroker: kafkaBroker, OtelExporterOtlpEndpoint: otelExporterOtlpEndpoint, OtelServiceName: otelServiceName}, nil
}
