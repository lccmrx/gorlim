package gorlim

import "time"

type timeframe int

const (
	RPS timeframe = iota
	RPM
)

var timeframeDuration = map[timeframe]time.Duration{
	RPS: time.Second,
	RPM: time.Minute,
}

var timeframeString = map[string]timeframe{
	"RPS": RPS,
	"RPM": RPM,
}

func WithTimeframe(timeframe timeframe) optionFunc {
	return func(r *RateLimiter) {
		r.timeframe = timeframe
	}
}
