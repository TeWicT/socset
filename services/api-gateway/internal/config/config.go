package config

type Config struct {
	HttpAddr     string
	GRPCAddrAuth string
}

func CreateConfig(httpAddr, grpcAddrAuth string) *Config {
	return &Config{HttpAddr: httpAddr, GRPCAddrAuth: grpcAddrAuth}
}
