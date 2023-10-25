/*
Package queue provides a Pure Golang thread-safe and unlimited-size generics in-memory message queue implementation
that supports async enqueue and blocking dequeue.
It's alternative way to communicate between goroutines compared to `channel`.
The implementation of this in-memory message queue uses sync.Cond instead of channel.

This queue is low-level and simple library, it's not a full-featured message queue.
If your use-case requires a limited-size queue and blocking enqueue, please use a channel instead.
For advanced use-cases like distributed queue, persistent message please use a message broker like Kafka, RabbitMQ, NATES or NSQ instead.

You can build any advanced message queue on top of this queue (use this queue for under the hood)
like an advance message queue like a single-producer with multiple-consumers queue,
broadcast system, multiple topics queue or any other use-cases.

# Basic usage

	q := New[string]()
	q.Enqueue("Foo")

	item, ok := q.Dequeue()
	if !ok {
	  panic("item should be dequeued")
	}
	fmt.Println(item) // Output: Foo

Queue is unlimited capacity, so you can enqueue as many as you want without blocking or dequeue required:

	q := New[string]()
	for i := 0; i < 1000000; i++ {
	  q.Enqueue("Foo")
	}

# Concurrency usage

		q := New[string]()

		go func() {
		  for {
		    item, ok := q.Dequeue()
		    if !ok {
		      break
		    }
		    fmt.Println(item)
		  }
		}()

		data := []string{"Foo", "Bar", "Baz"}
		for _, item := range data {
		  q.Enqueue(item)
		}

	  q.Close() // close queue to stop goroutine
*/
package queue
