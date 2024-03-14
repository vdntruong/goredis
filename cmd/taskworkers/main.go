package main

import (
	"log"
	"time"

	"goredis/tasks"

	"github.com/hibiken/asynq"
)

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     "localhost:6379",
		Password: "mypassword",
	})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            on background side, testing for Concurrency number 10 goroutine at the time.
	// ------------------------------------------------------
	// for i := 10; i < 30; i++ {
	// 	task, err := tasks.NewEmailDeliveryTask(i, fmt.Sprintf("some:template:%s", i))
	// 	if err != nil {
	// 		log.Fatalf("could not create task: %v", err)
	// 	}

	// 	info, err := client.Enqueue(task)
	// 	if err != nil {
	// 		log.Fatalf("could not enqueue task: %v", err)
	// 	}
	// 	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	// }
	// return

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------
	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------
	info, err = client.Enqueue(task, asynq.ProcessIn(24*time.Hour))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	task, err = tasks.NewImageResizeTask("https://example.com/myassets/example-image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
