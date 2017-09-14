package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := stan.Connect(
		"test-cluster",
		"test-client",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		return err
	}
	defer logCloser(conn)

	wg := &sync.WaitGroup{}

	var lastProcessed uint64
	var i int

	sub, err := conn.Subscribe("counter", func(msg *stan.Msg) {
		var processed bool

		if msg.Sequence > lastProcessed {
			processed = true
			atomic.SwapUint64(&lastProcessed, msg.Sequence)
		}

		// Add jitter..
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		i++

		var acked bool
		if i <= 5 {
			msg.Ack()
			// Mark it is done.
			wg.Done()
			acked = true
		} else if i == 9 {
			i = -5
		}

		// Print the value and whether it was redelivered.
		fmt.Printf("seq = %d [redelivered = %v, acked = %v, processed = %v]\n", msg.Sequence, msg.Redelivered, acked, processed)

	}, stan.SetManualAckMode(), stan.AckWait(time.Second))
	if err != nil {
		return err
	}
	defer logCloser(sub)

	// Publish up to 10.
	for i := 0; i < 10; i++ {
		wg.Add(1)

		err := conn.Publish("counter", nil)
		if err != nil {
			return err
		}
	}

	// Wait until all messages have been processed.
	wg.Wait()

	return nil
}
