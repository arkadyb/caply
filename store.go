package ratelimiter

import "time"

// Store defines behaviour of rate limiter store object
type Store interface {
	// Get returns number of processed operations by given key
	Get(key string) (int, error)
	// AddOne increases by one a number of served operations
	// ExpiresIn indicates for how long given operation timout should last
	AddOne(key string, expiresIn time.Duration) error
}
