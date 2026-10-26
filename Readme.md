# Go Rate Limiter

A rate-limiting service built with **Go (Golang)** using the **Token Bucket algorithm**, HTTP middleware, in-memory storage, and Redis. It limits the number of requests a client can make within a given time period.

## Features

- **Token Bucket Algorithm** — controls request rates using a configurable token refill rate and bucket capacity.
- **IP-Based Rate Limiting** — maintains rate limits independently for each client IP.
- **HTTP Middleware** — integrates rate limiting into HTTP request handling.
- **In-Memory Storage** — stores client buckets using Go maps and mutexes.
- **Redis Storage** — supports shared rate-limit state across application instances.
- **Atomic Redis Operations** — uses Lua scripting to refill and consume tokens atomically.
- **Error Handling** — handles invalid configuration and storage errors.
- **Automatic Redis Key Expiration** — removes inactive Redis rate-limit keys using expiration.
- **Unit and Integration Tests** — tests configuration validation, rate limiting, middleware, and Redis behavior.

## Tech Stack

- Go
- `net/http`
- Redis
- `github.com/redis/go-redis/v9`
- Lua scripting
- Go testing package

## How It Works

The project uses the Token Bucket algorithm.

Each client has a bucket containing tokens. Every incoming request consumes one token if a token is available. Tokens are replenished over time at a configured rate, up to the bucket's maximum capacity.

When no tokens remain, the request is rejected.

### Example Configuration

```go
limiter, err := NewRateLimiter(store, 2, 5)
if err != nil {
    log.Fatal(err)
}
```

In this example:

- Refill rate: 2 tokens per second.
- Bucket capacity: 5 tokens.
- Initial burst: up to 5 requests, assuming the bucket starts full.

The exact sustained request rate depends on the refill rate and request timing.

### Request Flow

1. A client sends an HTTP request.
2. Middleware extracts the client's IP address.
3. The rate limiter checks the client's token bucket.
4. If a token is available, the request continues to the handler.
5. If no token is available, the server returns `429 Too Many Requests`.
6. If the storage operation fails, the middleware returns `503 Service Unavailable`.

## Project Structure

```text
rate_limiter/
├── main.go
├── limiter.go
├── middleware.go
├── store.go
├── memory_store.go
├── redis_store.go
├── limiter_test.go
├── limiter_config_test.go
├── middleware_test.go
├── redis_store_test.go
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go installed on your machine.
- Redis installed and running if you want to use Redis storage.

Check your Go installation:

```powershell
go version
```

## Installation

Clone the repository:

```bash
git clone <your-repository-url>
cd rate_limiter
```

Download dependencies:

```bash
go mod tidy
```

## Running with Redis

Start your local Redis server and ensure it is accessible at `localhost:6379`, or update the Redis address in your configuration.

Run the application:

```bash
go run .
```

The server starts on:

```text
http://localhost:8080
```

If Redis is required by your current `main.go`, the application will fail to start if it cannot connect.

## Testing

Run all tests:

```bash
go test -v ./...
```

The test suite covers:

- Invalid rate-limiter configuration.
- Token consumption and request rejection.
- Independent limits for different clients.
- HTTP middleware responses.
- Redis token-bucket behavior.

The Redis integration test may be skipped if Redis is unavailable, depending on the test implementation.

## Testing the API

Send requests to the server:

```bash
curl -i http://localhost:8080/
```

Or use PowerShell:

```powershell
Invoke-WebRequest http://localhost:8080/
```

When the request is within the configured limit, the handler should respond successfully. Once the client's available tokens are exhausted, the server should return HTTP `429`.

To test consistently, use the same client and avoid waiting long enough for tokens to refill.

## HTTP Status Codes

| Status | Meaning |
|---|---|
| `200 OK` | Request accepted and handled successfully |
| `400 Bad Request` | Invalid client address |
| `429 Too Many Requests` | Client exceeded the rate limit |
| `503 Service Unavailable` | Rate-limit storage operation failed |

## Key Concepts Demonstrated

This project provides practical experience with:

- Interfaces and dependency injection.
- Structs and methods.
- Mutexes and concurrency safety.
- HTTP handlers and middleware.
- Token Bucket rate-limiting algorithms.
- Redis data structures and Lua scripting.
- Atomic operations in distributed systems.
- Error handling and configuration validation.
- Unit testing and integration testing.

## Future Improvements

- Graceful HTTP server shutdown.
- HTTP read, write, and idle timeouts.
- Configurable rate limits through environment variables.
- Docker and Docker Compose support.
- Forwarded-IP handling for trusted reverse proxies.
- Metrics and observability.
- Benchmarks and concurrent-load testing.
- More robust Redis failure handling.

## Learning Outcome

This project demonstrates how to build a reusable rate limiter in Go, integrate it into an HTTP server, and support both local in-memory state and shared Redis-backed state. It provides a foundation for exploring production backend concerns such as concurrency, distributed coordination, reliability, and testing.

## License

Choose a license before publishing this project publicly. For example, the MIT License is a common option for open-source projects.
