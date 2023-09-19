package config

type (
	Config struct {
		HttpServerAddr     string `env:"HTTP_SERVER_ADDR" validate:"required"`
		CollectorGrpcAddr  string `env:"COLLECTOR_GRPC_ADDR" validate:"required"`
		ManagementGrpcAddr string `env:"MANAGEMENT_GRPC_ADDR" validate:"required"`
	}
)
