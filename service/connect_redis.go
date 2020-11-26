package service

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func ConnectRedisPool() redis.Conn {
	connPool := &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow:    nil,
		MaxIdle:         1,
		MaxActive:       10,
		IdleTimeout:     180 * time.Second,
		Wait:            true,
		MaxConnLifetime: 0,
	}

	return connPool.Get()
}

