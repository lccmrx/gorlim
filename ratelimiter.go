package gorlim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const (
	RATE_LIMITER_MAX_REQUESTS_PER_TIME = "RATE_LIMITER_MAX_REQUESTS_PER_TIME"
	RATE_LIMITER_TIMEFRAME             = "RATE_LIMITER_TIMEFRAME"
	RATE_LIMITER_HEADER_LIMITER        = "RATE_LIMITER_HEADER_LIMITER"
)

type RateLimiter struct {
	backend Backend

	maxRequestsPerTime int

	timeframe timeframe

	enabledHeaderLimiter bool
	headerMap            map[string]int
}

type optionFunc func(*RateLimiter)

func New(backend Backend, opts ...optionFunc) *RateLimiter {
	ratelimiter := new(RateLimiter)
	ratelimiter.backend = backend

	ratelimiter.headerMap = make(map[string]int)
	ratelimiter.maxRequestsPerTime = 30

	ratelimiter.setFromEnv()

	for _, opt := range opts {
		opt(ratelimiter)
	}
	return ratelimiter
}

func WithHeaderLimit(header string, maxRequestsPerTime int) optionFunc {
	return func(r *RateLimiter) {
		r.enabledHeaderLimiter = true

		r.headerMap[header] = maxRequestsPerTime
	}
}

func WithRequestLimit(maxRequestsPerTime int) optionFunc {
	return func(r *RateLimiter) {
		r.maxRequestsPerTime = maxRequestsPerTime
	}
}

func Wrap(ratelimiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// apply rate limit logic
			err := ratelimiter.execute(w, r)
			if err != nil {
				json.NewEncoder(w).Encode(struct {
					Error string `json:"error"`
				}{
					Error: err.Error(),
				})
				return
			}

			// forward request to handler
			next.ServeHTTP(w, r)
		},
	)
}

func (ratelimiter *RateLimiter) setFromEnv() {
	maxRequestsPerTime := os.Getenv(RATE_LIMITER_MAX_REQUESTS_PER_TIME)
	if maxRequestsPerTime != "" {
		ratelimiter.maxRequestsPerTime, _ = strconv.Atoi(maxRequestsPerTime)
	}

	timeframe := os.Getenv(RATE_LIMITER_TIMEFRAME)
	if timeframe != "" {
		if timeframe, ok := timeframeString[timeframe]; ok {
			ratelimiter.timeframe = timeframe

		}
	}

	headerMap := os.Getenv(RATE_LIMITER_HEADER_LIMITER)
	if headerMap != "" {
		ratelimiter.enabledHeaderLimiter = true

		json.Unmarshal([]byte(headerMap), &ratelimiter.headerMap)
	}
}

func (ratelimiter *RateLimiter) execute(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	limiterKey := getRealIP(r)
	limiterDuration := timeframeDuration[ratelimiter.timeframe]
	limiterMaxRequests := ratelimiter.maxRequestsPerTime

	if ratelimiter.enabledHeaderLimiter {
		for header, maxRequestsPerTime := range ratelimiter.headerMap {
			if header := r.Header.Get(header); header != "" {
				limiterKey = header
				limiterMaxRequests = maxRequestsPerTime
			}
		}
	}

	limiterScore, err := ratelimiter.backend.GetScore(ctx, limiterKey, limiterDuration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to get backend key score: %s", err)
	}

	if limiterMaxRequests <= limiterScore {
		w.WriteHeader(http.StatusTooManyRequests)
		return fmt.Errorf("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	ratelimiter.backend.IncreaseScore(ctx, limiterKey, limiterDuration)
	return nil
}

func getRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
