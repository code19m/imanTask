package ports

import (
	"collector/internal/domain"
	"context"
)

type Service interface {
	ScheduleCollectPosts(ctx context.Context, payload domain.CollectPostPayload) (taskID string, err error)
	CollectPosts(ctx context.Context, startPage int) error
	GetTaskStateByID(ctx context.Context, taskID string) (string, error)
}
