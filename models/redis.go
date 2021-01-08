package models

import (
	"github.com/beego/beego/core/config"
	"github.com/garyburd/redigo/redis"
)

func ConnectRedisPool() redis.Conn {
	connPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", RedisLink())
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	return connPool.Get()
}

func ConnectRedis() redis.Conn {
	conn, _ := redis.Dial("tcp", RedisLink())
	return conn
}

func RedisLink() string {
	redis_, _ := config.String("redis_addr")
	return redis_
}