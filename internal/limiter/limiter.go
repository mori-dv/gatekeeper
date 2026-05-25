package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	redis  *redis.Client
	limit  int
	window time.Duration
}

func New(
	redis *redis.Client,
	limit int,
	window time.Duration,
) *Limiter {
	return &Limiter{
		redis:  redis,
		limit:  limit,
		window: window,
	}
}

func (l *Limiter) Allow(
	ctx context.Context,
	key string,
) (bool, int, error) {
	count, err := l.redis.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	if count == 1 {
		err := l.redis.Expire(ctx, key, l.window).Err()
		if err != nil {
			return false, 0, err
		}
	}

	remaining := l.limit - int(count)

	if remaining < 0 {
		remaining = 0
	}

	allowed := count <= int64(l.limit)

	return allowed, remaining, nil
}

func (l *Limiter) GetLimit() int {
	return l.limit
}