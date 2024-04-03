package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
		return fmt.Errorf("failed to marsal task payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailDelivery, jsonPayload)
	if err != nil {
		return fmt.Errorf("failed to create new task from task payload: %w", err)
	}

	info, err := r.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"enqueued task type=%s queue=%s retried=%d id=%s payload=%s\n",
		info.Type, info.Queue, info.Retried, info.ID, string(info.Payload),
	)
	return nil
}

func (r *RedisDistributor) DistributeImageTask(ctx context.Context, payload ImageResizePayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marsal task payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailDelivery, jsonPayload)
	if err != nil {
		return fmt.Errorf("failed to create new task from task payload: %w", err)
	}

	info, err := r.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"enqueued task type=%s queue=%s retried=%d id=%s payload=%s\n",
		info.Type, info.Queue, info.Retried, info.ID, string(info.Payload),
	)
	return nil
}

func (r *RedisDistributor) Close() error {
	return r.client.Close()
}
