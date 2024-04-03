package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "very-critical"
	QueueDefault  = "just-normal"
)

type EmailTaskProcessor interface {
	ProcessEmailTask(ctx context.Context, task *asynq.Task) error
}

type ImageTaskProcessor interface {
	ProcessImageTask(ctx context.Context, task *asynq.Task) error
}

type TaskProcessor struct {
	server *asynq.Server
}

func NewTaskProcessor(server *asynq.Server) *TaskProcessor {
	return &TaskProcessor{server: server}
}

func (t *TaskProcessor) ProcessImageTask(ctx context.Context, task *asynq.Task) error {
	var payload ImageResizePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarsal image task: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf(
		"processed image task type=%s, image_source=%s, payload=%s\n",
		task.Type(), payload.SourceURL, string(task.Payload()),
	)
	return nil
}

func (t *TaskProcessor) ProcessEmailTask(ctx context.Context, task *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarsal email task: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf(
		"processed email task type=%s, user_id=%d, payload=%s\n",
		task.Type(), payload.UserID, string(task.Payload()),
	)
	return nil
}

func (t *TaskProcessor) Start() error {
	muxSvr := asynq.NewServeMux()
	muxSvr.HandleFunc(TypeEmailDelivery, t.ProcessEmailTask)
	muxSvr.HandleFunc(TypeImageResize, t.ProcessImageTask)
	return t.server.Start(muxSvr)
}

func (t *TaskProcessor) Stop() {
	t.server.Stop()
	t.server.Shutdown()
}
