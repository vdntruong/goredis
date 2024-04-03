package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type Distributor interface {
	DistributeEmailTask(ctx context.Context, payload EmailDeliveryPayload) error
	DistributeImageTask(ctx context.Context, payload ImageResizePayload) error
}

type RedisDistributor struct {
	client *asynq.Client
}

func NewRedisDistributor(client *asynq.Client) *RedisDistributor {
	return &RedisDistributor{client: client}
}

func (r *RedisDistributor) DistributeEmailTask(ctx context.Context, payload EmailDeliveryPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marsal email task payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailDelivery, jsonPayload)
	if err != nil {
		return fmt.Errorf("failed to create new task from payload: %w", err)
	}

	opt := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(10 * time.Second), // process the task after 10s ~ delay
		asynq.Queue(QueueCritical),
	}
	info, err := r.client.Enqueue(task, opt...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"enqueued email task type=%s queue=%s retried=%d id=%s payload=%s\n",
		info.Type, info.Queue, info.Retried, info.ID, string(info.Payload),
	)
	return nil
}

func (r *RedisDistributor) DistributeImageTask(ctx context.Context, payload ImageResizePayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marsal image task payload: %w", err)
	}

	task := asynq.NewTask(TypeImageResize, jsonPayload)
	if err != nil {
		return fmt.Errorf("failed to create new task from payload: %w", err)
	}

	opt := []asynq.Option{
		asynq.MaxRetry(0),
		asynq.Queue(QueueDefault),
		asynq.Timeout(2 * time.Second), // the processor's context with go with a deadline
	}
	info, err := r.client.Enqueue(task, opt...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"enqueued image task type=%s queue=%s retried=%d id=%s payload=%s\n",
		info.Type, info.Queue, info.Retried, info.ID, string(info.Payload),
	)
	return nil
}

func (r *RedisDistributor) Close() error {
	return r.client.Close()
}
