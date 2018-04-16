package event

import (
	//	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	var ch = make(chan struct{}, 1)
	var name = `event_name`
	var expect = 123
	On(name, func(raw []byte) {
		var value int
		Parse(raw, &value)
		if value != expect {
			t.Errorf(`expect: %s; got: %s`, expect, value)
		}
		ch <- struct{}{}
	})
	Trigger(name, 123)
	<-ch
}
