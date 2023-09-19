package grpcsrv

import (
	"collector/gen/pb"
	"collector/internal/ports"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	*pb.UnimplementedTaskServiceServer
	grpcSrv *grpc.Server
	service ports.Service
}

func New(service ports.Service) *server {
	return &server{
		UnimplementedTaskServiceServer: &pb.UnimplementedTaskServiceServer{},
		service:                        service,
	}
}

func (s *server) Run(addr string) {
	gprcLogger := grpc.UnaryInterceptor(GrpcLogger)

	s.grpcSrv = grpc.NewServer(gprcLogger)

	pb.RegisterTaskServiceServer(s.grpcSrv, s)

	reflection.Register(s.grpcSrv)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = s.grpcSrv.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func (s *server) ShutDown() {
	s.grpcSrv.GracefulStop()
}
