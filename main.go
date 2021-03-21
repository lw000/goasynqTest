package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"goasyncnq/tasks"
	"log"
	"time"
)

const redisAddr = "127.0.0.1:6379"

func TestClientMain() {
	r := asynq.RedisClientOpt{Addr: redisAddr}
	c := asynq.NewClient(r)
	defer c.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	t := tasks.NewEmailDeliveryTask(42, "some:template:id")
	res, err := c.Enqueue(t)
	if err != nil {
		log.Fatal("could not enqueue task: %v", err)
	}
	fmt.Printf("Enqueued Result: %+v\n", res)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	t = tasks.NewEmailDeliveryTask(42, "other:template:id")
	res, err = c.Enqueue(t, asynq.ProcessIn(24*time.Hour))
	if err != nil {
		log.Fatal("could not schedule task: %v", err)
	}
	fmt.Printf("Enqueued Result: %+v\n", res)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	c.SetDefaultOptions(tasks.TypeImageResize, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))

	t = tasks.NewImageResizeTask("some/blobstore/path")
	res, err = c.Enqueue(t)
	if err != nil {
		log.Fatal("could not enqueue task: %v", err)
	}
	fmt.Printf("Enqueued Result: %+v\n", res)

	// ---------------------------------------------------------------------------
	// Example 4: Pass options to tune task processing behavior at enqueue time.
	//            Options passed at enqueue time override default ones, if any.
	// ---------------------------------------------------------------------------

	t = tasks.NewImageResizeTask("some/blobstore/path")
	res, err = c.Enqueue(t, asynq.Queue("critical"), asynq.Timeout(30*time.Second))
	if err != nil {
		log.Fatal("could not enqueue task: %v", err)
	}
	fmt.Printf("Enqueued Result: %+v\n", res)
}

func TestServerMain() {
	r := asynq.RedisClientOpt{Addr: redisAddr}

	srv := asynq.NewServer(r, asynq.Config{
		// Specify how many concurrent workers to use
		Concurrency: 10,
		// Optionally specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		// See the godoc for other configuration options
	})

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func main() {
	time.AfterFunc(time.Second*2, func() {
		TestClientMain()
	})

	TestServerMain()
}
