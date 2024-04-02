package tasks

import (
	"context"
	"errors"
	"log"
	"math/rand"
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
	defer func() {
		wg.Done()
	}()

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

func background(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	muxSvr := asynq.NewServeMux()
	muxSvr.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	muxSvr.Handle(TypeImageResize, NewImageProcessor())

	if err := srv.Start(muxSvr); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		s := <-sigs
		if s == unix.SIGTSTP {
			log.Println("Shutdown Server ...")
			srv.Stop() // stop processing new tasks
			continue
		}
		break
	}

	srv.Shutdown()
	log.Println("Background Server exiting")
}

func client() {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPass,
	})
	defer func() { _ = client.Close() }()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------
	task, err := NewEmailDeliveryTask(rand.Int(), "some:template:01")
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
	task02, err := NewEmailDeliveryTask(rand.Int(), "some:template:02")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("could not load location Shanghai: %v", err)
	}
	at := time.Date(2024, time.April, 2, 17, 28, 0, 0, loc)

	info, err = client.Enqueue(task02, asynq.ProcessAt(at))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------
	task, err = NewImageResizeTask("https://example.com/myassets/example-image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func Run() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		serve(&wg)
	}()

	wg.Add(1)
	go func() {
		background(&wg)
	}()

	client()
	wg.Wait()
}
