package config

type (
	Config struct {
		GrpcServerAddr string `env:"GRPC_SERVER_ADDR" validate:"required"`
	}
)
