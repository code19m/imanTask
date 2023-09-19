package service

import (
	"collector/internal/domain"
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
)

func (s *service) ScheduleCollectPosts(ctx context.Context, payload domain.CollectPostPayload) (taskID string, err error) {
	taskPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	task := asynq.NewTask(domain.CollectPostsTaskName, taskPayload, asynq.MaxRetry(0))

	info, err := s.taskClient.EnqueueContext(ctx, task)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}

func (s *service) GetTaskStateByID(ctx context.Context, taskID string) (string, error) {
	info, err := s.inspector.GetTaskInfo("default", taskID)

	if errors.Is(err, asynq.ErrTaskNotFound) {
		return "not found or processed", nil
	}

	if err != nil {
		return "", err
	}

	return info.State.String(), nil
}
