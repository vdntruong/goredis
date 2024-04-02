package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code
	_, cancelFunc := context.WithTimeout(ctx, time.Second*3)
	defer cancelFunc()
	time.Sleep(2 * time.Second)
	log.Printf("Sent Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)

	return nil
}
