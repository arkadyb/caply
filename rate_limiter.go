package ratelimiter

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type RateLimiter interface {
	LimitExceeded(commandName string) (bool, error)
}

func NewFixedTimeWindowRateLimiter(maxRequests int, perPeriod time.Duration, store Store) (*FixedTimeWindowRateLimiter, error) {
	if perPeriod < 1*time.Second || perPeriod > 1*time.Hour {
		return nil, errors.New("perPeriod has to be between 1 second and 1 hour")
	}

	return &FixedTimeWindowRateLimiter{
		store,
		maxRequests,
		perPeriod,
	}, nil
}

type FixedTimeWindowRateLimiter struct {
	store       Store
	maxRequests int
	perPeriod   time.Duration
}

func (rt *FixedTimeWindowRateLimiter) LimitExceeded(commandName string) (bool, error) {
	bucketTimeStamp := 0
	now := time.Now()
	if rt.perPeriod < 1*time.Minute {
		bucketTimeStamp = now.Minute()
	} else if rt.perPeriod < 1*time.Hour {
		bucketTimeStamp = now.Hour()
	}

	key := fmt.Sprintf("%s_%d", commandName, bucketTimeStamp)
	val, err := rt.store.Get(key)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get current limit state for key %s", key)
	}
	if val >= rt.maxRequests {
		return true, nil
	} else {
		err = rt.store.AddOne(key, rt.perPeriod)
		if err != nil {
			return false, errors.Wrapf(err, "failed to increase rate limit for key %s", key)
		}
	}

	return false, nil
}
