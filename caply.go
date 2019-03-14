package caply

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

// NewCaply returns new instance of Caply
func NewCaply(maxOps int, perPeriod time.Duration, store Store) (*Caply, error) {
	if perPeriod < 1*time.Second || perPeriod > 1*time.Hour {
		return nil, errors.New("perPeriod has to be between 1 second and 1 hour")
	}

	return &Caply{
		store,
		maxOps,
		perPeriod,
	}, nil
}

// Caply is fixed window rate limiter where every next request is placed into respective time window batch
type Caply struct {
	store       Store
	maxRequests int
	perPeriod   time.Duration
}

// Exceeded finds specified operation in current time window and returns true if limit has already exceeded (false otherwise)
func (rt *Caply) Exceeded(opName string) (bool, error) {
	bucketTimeStamp := 0
	now := time.Now()
	if rt.perPeriod < 1*time.Minute {
		bucketTimeStamp = now.Minute()
	} else if rt.perPeriod < 1*time.Hour {
		bucketTimeStamp = now.Hour()
	}

	key := fmt.Sprintf("%s_%d", opName, bucketTimeStamp)
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
