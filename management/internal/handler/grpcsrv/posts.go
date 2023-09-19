package grpcsrv

import (
	"context"
	"management/gen/pb"
	"management/internal/domain"
	"management/pkg/mapper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetPosts(ctx context.Context, in *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	posts, count, err := s.service.GetPosts(ctx, in.GetOffset(), in.GetLimit())
	if err != nil {
		return nil, err
	}

	resp := new(pb.GetPostsResponse)
	resp.Count = count
	err = mapper.Map(posts, &resp.Posts)
	return resp, err
}

func (s *server) GetPostById(ctx context.Context, in *pb.GetPostByIdRequest) (*pb.GetPostByIdResponse, error) {
	post, err := s.service.GetPost(ctx, in.GetId())
	if err == domain.ErrPostNotFound {
		return nil, status.Errorf(codes.NotFound, "post with ID '%d' not found", in.GetId())
	}

	if err != nil {
		return nil, err
	}

	resp := new(pb.GetPostByIdResponse)
	err = mapper.Map(post, &resp.Post)
	return resp, err
}

func (s *server) UpdatePost(ctx context.Context, in *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id: required field")
	}
	if in.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id: required field")
	}
	if in.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title: required field")
	}
	if in.GetBody() == "" {
		return nil, status.Error(codes.InvalidArgument, "body: required field")
	}

	err := s.service.UpdatePost(ctx, domain.Post{
		Id:     int(in.GetId()),
		UserID: int(in.GetUserId()),
		Title:  in.GetTitle(),
		Body:   in.GetBody(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePostResponse{Success: true}, nil
}

func (s *server) DeletePost(ctx context.Context, in *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	err := s.service.DeletePost(ctx, in.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePostResponse{Success: true}, nil
}
