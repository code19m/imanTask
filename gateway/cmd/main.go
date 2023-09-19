package main

import (
	"context"
	"gateway/config"
	"gateway/gen/pb"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb.RegisterTaskServiceHandlerFromEndpoint(ctx, mux, cfg.CollectorGrpcAddr, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register collector endpoints")
	}

	err = pb.RegisterPostServiceHandlerFromEndpoint(ctx, mux, cfg.ManagementGrpcAddr, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register management endpoints")
	}

	if err := http.ListenAndServe(cfg.HttpServerAddr, mux); err != nil {
		log.Fatal().Err(err).Msg("failed to start http server")
	}
}
