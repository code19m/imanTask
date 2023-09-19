package grpcsrv

import (
	"collector/gen/pb"
	"collector/internal/domain"
	"context"
)

func (s *server) ScheduleCollectPosts(ctx context.Context, in *pb.ScheduleCollectPostsRequest) (*pb.ScheduleCollectPostsResponse, error) {
	payload := domain.CollectPostPayload{
		StartPage: in.StartPage,
	}

	taskID, err := s.service.ScheduleCollectPosts(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &pb.ScheduleCollectPostsResponse{Id: taskID}, nil
}
