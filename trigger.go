package event

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

func Trigger(channel string, data interface{}) {
	var bytes []byte
	var err error
	if data != nil {
		bytes, err = json.Marshal(data)
		if err != nil {
			panic(err)
		}
	}
	redisDo(func(conn redis.Conn) {
		if _, err := redis.Int(conn.Do(`publish`, channel, bytes)); err != nil {
			panic(err)
		}
	})
}
