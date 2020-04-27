package tcp

import (
	"sync"
	"time"
)

type ClientId int64

func ClientIdEmpty() ClientId {
	return ClientId(0)
}

func GenClientId() ClientId {
	return ClientId(time.Now().UnixNano())
}

func IntToClientId(id int) ClientId {
	return ClientId(id)
}

func Int64ToClientId(id int64) ClientId {
	return ClientId(id)
}

// serverID => Connect
type ConnectionsMap struct {
	sync.Map
}

func (c *ConnectionsMap) RLoad(key string) (value *Connect, ok bool) {
	var v interface{}
	v, ok = c.Load(key)
	if ok {
		value = v.(*Connect)
	} else {
		value = nil
	}
	return
}

