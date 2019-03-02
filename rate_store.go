package ratelimiter

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

func NewRedisStore(redisPool *redis.Pool) *RedisStore {
	return &RedisStore{
		redisPool,
	}
}

type RedisStore struct {
	redisPool *redis.Pool
}

func (rd *RedisStore) Get(id string) (int, error) {
	conn := rd.redisPool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("GET", id))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		}
		return 0, errors.Wrapf(err, "failed to get value by id %s", id)
	}

	return val, nil
}

func (rd *RedisStore) AddOne(id string, expiresIn time.Duration) error {
	conn := rd.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("MULTI")
	if err != nil {
		return errors.Wrapf(err, "failed add one for id %s", id)
	}
	_, err = conn.Do("INCR", id)
	if err != nil {
		return errors.Wrapf(err, "failed add one for id %s", id)
	}

	_, err = conn.Do("EXPIRE", id, expiresIn.Seconds())
	if err != nil {
		return errors.Wrapf(err, "failed add one for id %s", id)
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return errors.Wrapf(err, "failed add one for id %s", id)
	}

	return nil
}
