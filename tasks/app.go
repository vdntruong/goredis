package tasks

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"golang.org/x/sys/unix"
)

const redisAddr = "127.0.0.1:6379"
const redisPass = "mypassword"

func serve(wg *sync.WaitGroup) {
	defer wg.Done()

	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring",
		RedisConnOpt: asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass},
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal for HTTP Server")

	toutCtx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancelFunc()
	if err := srv.Shutdown(toutCtx); err != nil {
		log.Printf("HTTP Server shutdown returned an err: %v\n", err)
	}
	log.Println("HTTP Server exiting")
}

func processTasks(wg *sync.WaitGroup) {
	defer wg.Done()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				QueueCritical: 6,
				QueueDefault:  3,
			},
		},
	)

	processor := NewTaskProcessor(srv)
	if err := processor.Start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		s := <-sigs
		if s == unix.SIGTSTP {
			log.Println("Shutting down the server...")
			processor.Stop()
			continue
		}
		break
	}
	log.Println("Processor exited")
}

func distributeTasks(wg *sync.WaitGroup) {
	defer wg.Done()

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPass,
	})

	distributor := NewRedisDistributor(client)
	defer func() { _ = distributor.Close() }()

	for i := 0; i < 5; i++ {
		if err := distributor.DistributeEmailTask(
			context.TODO(),
			EmailDeliveryPayload{
				UserID:     i,
				TemplateID: fmt.Sprintf("email:template:%d", i),
			},
		); err != nil {
			log.Println("distribute email task err", err)
		}
	}

	for i := 0; i < 3; i++ {
		if err := distributor.DistributeImageTask(
			context.TODO(),
			ImageResizePayload{
				SourceURL: fmt.Sprintf("https://example..com/example-image%d.jpg", i),
			},
		); err != nil {
			log.Println("distribute email task err", err)
		}
	}

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------
	//task, err := NewEmailDeliveryTask(rand.Int(), "some:template:01")
	//if err != nil {
	//	log.Fatalf("could not create task: %v", err)
	//}
	//
	//info, err := client.Enqueue(task)
	//if err != nil {
	//	log.Fatalf("could not enqueue task: %v", err)
	//}
	//log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------
	//task02, err := NewEmailDeliveryTask(rand.Int(), "some:template:02")
	//if err != nil {
	//	log.Fatalf("could not create task: %v", err)
	//}
	//
	//loc, err := time.LoadLocation("Asia/Shanghai")
	//if err != nil {
	//	log.Fatalf("could not load location Shanghai: %v", err)
	//}
	//at := time.Date(2024, time.April, 2, 17, 28, 0, 0, loc)
	//
	//info, err = client.Enqueue(task02, asynq.ProcessAt(at))
	//if err != nil {
	//	log.Fatalf("could not schedule task: %v", err)
	//}
	//log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------
	//task, err = NewImageResizeTask("https://example.com/myassets/example-image.jpg")
	//if err != nil {
	//	log.Fatalf("could not create task: %v", err)
	//}
	//
	//info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	//if err != nil {
	//	log.Fatalf("could not enqueue task: %v", err)
	//}
	//log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func Run() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		serve(&wg)
	}()

	wg.Add(1)
	go func() {
		processTasks(&wg)
	}()

	wg.Add(1)
	distributeTasks(&wg)

	wg.Wait()
}
