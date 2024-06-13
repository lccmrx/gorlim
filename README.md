# Gorlim

Gorlim is a Go Rate Limiter package, built using the standard Go HTTP package.

It's goal is to be a simple, right out of the box usage for many scenarios.

It can be quickly setup into your Server Multiplexer, like:
```go
ratelimiter := gorlim.New(
    backend.NewRedis("redis:6379"),
)

http.Handle("GET /", gorlim.Wrap(ratelimiter,
    http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ...
        }),
))
```

## Backend

Backends is a user-defined store system. One can define it's own adaptation
to fill Gorlim's usage. But be aware that most of the rating score for a desired key will be done by the user's implementation.

Gorlim provides a ready-to-use Redis Backend, providing only the connection address.

## Configuration
Gorlim provides some easy to use 
environment variable to set some of it's features.

```env
RATE_LIMITER_MAX_REQUESTS_PER_TIME=100
RATE_LIMITER_TIMEFRAME=RPM
RATE_LIMITER_HEADER_LIMITER='{"x-api-key": 10}'
```

Although, one could also set those configurarions in a programmatic way that override these environment variables.
```go
ratelimiter := gorlim.New(
    backend.NewRedis("redis:6379"),
    gorlim.WithRequestLimit(100),
    gorlim.WithTimeframe(gorlim.RPM),
    gorlim.WithHeaderLimit("x-api-key", 10),
)
```

Gorlim provides an IP-based rate limiting, such as Header-based.

Gorlim does it's job to find the user's IP, but if set, a Header will be taken in consideration more importantly.

So in this example below:
```go
ratelimiter := gorlim.New(
    backend.NewRedis("redis:6379"),
    gorlim.WithRequestLimit(100),
    gorlim.WithHeaderLimit("x-api-key", 10),
)
```

Even having both configurations, if the Header `x-api-key` is found on the incoming request, the IP request limit (100) won't matter to Gorlim.

Multiple Headers can be set, but the latter will precede on Gorlim's order.


## Colaboration

Feel free to reach me and collaborate with a simple tool for Go community!
