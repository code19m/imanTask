package worker

import (
	"collector/internal/domain"
	"collector/internal/ports"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type worker struct {
	server  *asynq.Server
	service ports.Service
}

func New(redisSrv *asynq.Server, service ports.Service) *worker {
	return &worker{server: redisSrv, service: service}
}

func (w *worker) Run() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(domain.CollectPostsTaskName, w.ProcessCollectPosts)

	err := w.server.Run(mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start worker")
	}
}

func (w *worker) ShutDown() {
	w.server.Shutdown()
}
