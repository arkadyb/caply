package ratelimiter

import "time"

type Store interface {
	Get(id string) (int, error)
	AddOne(id string, expiresIn time.Duration) error
}
