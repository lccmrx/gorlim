package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lccmrx/gorlim"
)

var _ gorlim.Backend = (*Redis)(nil)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *Redis) IncreaseScore(ctx context.Context, key string, timeframe time.Duration) error {
	now := time.Now().Unix()

	r.client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: uuid.New().String(),
	})

	minTimestamp := float64(now - int64(timeframe))
	r.client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", minTimestamp))

	return nil
}

func (r *Redis) GetScore(ctx context.Context, key string, timeframe time.Duration) (int, error) {
	now := time.Now()
	min := fmt.Sprintf("%d", now.Add(-timeframe).Unix())
	max := fmt.Sprintf("%d", now.Unix())

	result, err := r.client.ZCount(ctx, key, min, max).Result()
	return int(result), err
}
