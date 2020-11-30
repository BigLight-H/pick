package models

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

func ConnectRedisPool() redis.Conn {
	connPool := &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", beego.AppConfig.String("redis_addr"))
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

func ConnectRedis() redis.Conn {
	conn, _ := redis.Dial("tcp", beego.AppConfig.String("redis_addr"))
	return conn
}