package models

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

func ConnectRedisPool() redis.Conn {
	connPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", beego.AppConfig.String("redis_addr"))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	return connPool.Get()
}

func ConnectRedis() redis.Conn {
	conn, _ := redis.Dial("tcp", beego.AppConfig.String("redis_addr"))
	return conn
}