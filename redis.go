package event

import (
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

func getRedisUrl() string {
	switch os.Getenv(`GOENV`) {
	case `production`, `preview`:
		return `redis://:NabXadCcg0_ze08cSMFQ97JCGHKvAL6C@10.249.1.131:6379/0`
	case `qa`:
		return `redis://:@localhost:6379/0`
	default:
		return `redis://:@localhost:6379/0`
	}
}

var redisPool *redis.Pool

func getRedisPool() *redis.Pool {
	if redisPool == nil {
		redisPool = &redis.Pool{
			MaxIdle:     32,
			MaxActive:   32,
			IdleTimeout: 600 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					getRedisUrl(),
					redis.DialConnectTimeout(time.Second),
					redis.DialReadTimeout(time.Second),
					redis.DialWriteTimeout(time.Second),
				)
			},
		}
	}
	return redisPool
}

func redisDo(work func(redis.Conn)) {
	conn := getRedisPool().Get()
	defer conn.Close()
	work(conn)
}

func subscribeConn() (redis.Conn, error) {
	return redis.DialURL(
		getRedisUrl(),
		redis.DialConnectTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
	)
}
