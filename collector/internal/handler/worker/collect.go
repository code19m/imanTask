package worker

import (
	"collector/internal/domain"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

func (w *worker) ProcessCollectPosts(ctx context.Context, task *asynq.Task) error {
	var payload domain.CollectPostPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}
	return w.service.CollectPosts(ctx, int(payload.StartPage))
}
