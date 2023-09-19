package ports

import (
	"context"
	"management/internal/domain"
)

type Service interface {
	GetPosts(ctx context.Context, offset int32, limit int32) ([]domain.Post, int32, error)
	GetPost(ctx context.Context, id int32) (domain.Post, error)
	UpdatePost(ctx context.Context, payload domain.Post) error
	DeletePost(ctx context.Context, id int32) error
}
