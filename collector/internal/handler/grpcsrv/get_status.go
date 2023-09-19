package grpcsrv

import (
	"collector/gen/pb"
	"context"
)

func (s *server) GetTaskState(ctx context.Context, in *pb.GetTaskStateRequest) (*pb.GetTaskStateResponse, error) {
	status, err := s.service.GetTaskStateByID(ctx, in.TaskId)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskStateResponse{Status: status}, nil
}
