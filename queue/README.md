[![Go Reference](https://pkg.go.dev/badge/github.com/Cleverse/go-utilities/queue.svg)](https://pkg.go.dev/github.com/Cleverse/go-utilities/queue)
[![Report card](https://goreportcard.com/badge/github.com/Cleverse/go-utilities/queue)](https://goreportcard.com/report/github.com/Cleverse/go-utilities/queue)

# queue

Minimalist and zero-dependency thread-safe and unlimited-size generics in-memory message queue implementation
that supports async enqueue and blocking dequeue. \
It's alternative way to communicate between goroutines compared to `channel`

> **Note:** \
> This package is not intended to be used as a distributed message queue. For advanced use-cases like distributed queue, persistent message please use a message broker like Kafka, RabbitMQ, NATES or NSQ instead.
>
> And if your use-case requires a limited-size queue and blocking enqueue, please use a channel instead.

This package is low-level and simple queue library, it's not a full-featured message queue. \
You can build any advanced message queue on top of this queue (use this queue for under the hood)
like an advance message queue like a single-producer with multiple-consumers queue,
broadcast system, multiple topics queue or any other use-cases.

## Installation

```shell
go get github.com/Cleverse/go-utilities/queue
```
