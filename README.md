[![Go Reference](https://pkg.go.dev/badge/github.com/Cleverse/go-utilities.svg)](https://pkg.go.dev/github.com/Cleverse/go-utilities)

# go-utilities

Miscellaneous useful shared Go packages by [Cleverse](https://about.cleverse.com)

## utils

Minimalist pure Golang optimized generic utilities for Cleverse projects.

[See here](utils/README.md).

## errors

Minimalist and zero-dependency errors library with stacktrace support for Go (for wrapping and formatting an errors).

[See here](errors/README.md).

## nullable

A safe way to represent nullable primitive values in Go. Supports JSON serialization.

[See here](nullable/README.md).

## queue

Minimalist and zero-dependency low-level and simple queue library for thread-safe and unlimited-size generics in-memory message queue library for Go (async enqueue and blocking dequeue supports).\
The alternative way to communicate between goroutines compared to `channel`

[See here](queue/README.md).

## address

High efficient and minimal utilities library that will help you to work with Ethereum addresses easier. (a [go-ethereum](https://github.com/ethereum/go-ethereum) helper library)

[See here](address/README.md).

## fixedpoint

A [shopspring/decimal](https://github.com/shopspring/decimal) wrapper library for fixed point arithmetic operations in Cleverse projects.

[See here](fixedpoint/README.md).

## logger

A logger utility library, with support for TEXT and JSON logging (with optional GCP log format support). Supports embedding log attributes in context.

[See here](logger/README.md).

## httpclient

Simple [valyala/fasthttp](https://github.com/valyala/fasthttp) wrapper library with user-friendly interface and built-in request/response acquire and release.

[See here](httpclient/go.mod).

## postgres

`postgres` provides a wrapper around `pgx/v5` for connecting to PostgreSQL databases. It simplifies configuration, connection pooling, and integrates structured logging with `slog`.

[See here](postgres/README.md).

## redis

`redis` provides a wrapper around `go-redis` (v9) for connecting to Redis servers. It simplifies configuration and integrates structured logging with `slog`.

[See here](redis/README.md).

## cloudkms

`cloudkms` provides a simplified wrapper around the Google Cloud Key Management Service (KMS) API, making it easier to encrypt and decrypt data using Cloud KMS keys.

[See here](cloudkms/README.md).

## encryption

`encryption` provides simple and secure encryption utilities using ChaCha20-Poly1305, with base64 encoding/decoding support.

[See here](encryption/README.md).

## errs

`errs` defines a set of common application errors, built on top of `cockroachdb/errors`.

[See here](errs/README.md).

## automaxprocs

Automatically set GOMAXPROCS to match Linux container CPU quota by calling `automaxprocs.Init()`.

[See here](automaxprocs/README.md).
