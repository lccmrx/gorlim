package gorlim

import (
	"context"
	"time"
)

type Backend interface {
	GetScore(context.Context, string, time.Duration) (int, error)
	IncreaseScore(context.Context, string, time.Duration) error
}
