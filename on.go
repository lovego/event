package event

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

var events = struct {
	sync.RWMutex
	m map[string]func([]byte)
}{m: make(map[string]func([]byte))}

var subConn redis.PubSubConn

func On(channel string, handler func([]byte)) {
	if subConn.Conn == nil {
		subConn = getSubConn()
		startListen()
	}
	err := subConn.Subscribe(channel)
	if err != nil {
		log.Panic(err)
	}
	events.Lock()
	events.m[channel] = handler
	events.Unlock()
}

func getSubConn() redis.PubSubConn {
	conn, err := subscribeConn()
	if err != nil {
		log.Panic(err)
	}
	return redis.PubSubConn{Conn: conn}
}

func startListen() {
	go func() {
		for i := 0; i < 10000; i++ {
			func() {
				defer recover()
				listen(i)
			}()
		}
	}()
}

func listen(reconnect int) {
	defer subConn.Close()

	if reconnect > 0 {
		subConn = getSubConn()
		events.RLock()
		for channel, _ := range events.m {
			err := subConn.Subscribe(channel)
			if err != nil {
				log.Panic(err)
			}
		}
		events.RUnlock()
	}

Loop:
	for {
		switch n := subConn.Receive().(type) {
		case redis.Message:
			func() {
				defer recover()
				events.RLock()
				handler := events.m[n.Channel]
				events.RUnlock()
				handler(n.Data)
			}()
		case redis.Subscription:
			if reconnect != 0 {
				log.Printf("event Subscription(reconnect: %d): %s %s %d\n", reconnect, n.Kind, n.Channel, n.Count)
			}
		case error:
			log.Printf("event error: %v\n", n)
			break Loop
		}
	}
}

func Parse(raw []byte, v interface{}) {
	if raw == nil || len(raw) == 0 {
		return
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		log.Panic(err)
	}
	return
}
