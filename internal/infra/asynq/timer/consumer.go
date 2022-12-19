package timer

import (
	"fmt"
	"github.com/hibiken/asynq"
)

type Consumer struct {
	server    *asynq.Server
	processor asynq.Handler
}

func NewConsumer(server *asynq.Server, processor asynq.Handler) (*Consumer, error) {
	if server == nil {
		return nil, fmt.Errorf("asynq server is not set up")
	}

	if processor == nil {
		return nil, fmt.Errorf("processor is not set up")
	}

	return &Consumer{
		server:    server,
		processor: processor,
	}, nil
}

func (c *Consumer) Run() error {
	mux := asynq.NewServeMux()
	mux.Handle(TypeName, c.processor)
	return c.server.Run(mux)
}
