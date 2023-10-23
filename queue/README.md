[![GoDoc](https://godoc.org/github.com/Cleverse/go-utilities/queue?status.svg)](http://godoc.org/github.com/Cleverse/go-utilities/queue)
[![Report card](https://goreportcard.com/badge/github.com/Cleverse/go-utilities/queue)](https://goreportcard.com/report/github.com/Cleverse/go-utilities/queue)

# queue

A Pure Golang thread-safe and unlimited-size generics in-memory message queue implementation
that supports async enqueue and blocking dequeue. \
It's alternative way to communicate between goroutines compared to `channel`

> **Note:** \
> This package is not intended to be used as a distributed message queue. For advanced use-cases like distributed queue, persistent message please use a message broker like Kafka, RabbitMQ, NATES or NSQ instead.
>
> And if your use-case requires a limited-size queue and blocking enqueue, please use a channel instead.

## Installation

```shell
go get github.com/Cleverse/go-utilities/queue
```
