package config

type (
	Config struct {
		GrpcServerAddr string `env:"GRPC_SERVER_ADDR" validate:"required"`

		Redis redis
	}

	redis struct {
		Addr string `env:"REDIS_ADDR" validate:"required"`
		Pass string `env:"REDIS_PASSWORD"`
		DB   int    `env:"REDIS_DB"`
	}
)
