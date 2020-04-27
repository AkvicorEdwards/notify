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

type TeamId int64

func TeamIdEmpty() TeamId {
	return TeamId(0)
}

func GenTeamId() TeamId {
	return TeamId(time.Now().UnixNano())
}

func IntToTeamId(id int) TeamId {
	return TeamId(id)
}

func Int64ToTeamId(id int64) TeamId {
	return TeamId(id)
}

// UUID => *Connect
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

// api.clientId => *Client
type ConnMap struct {
	sync.Map
}

func (c *ConnMap) RLoad(key ClientId) (value *Client, ok bool) {
	var v interface{}
	v, ok = c.Load(key)
	if ok {
		value = v.(*Client)
	} else {
		value = nil
	}
	return
}

// TeamId => *Team
type TeamMap struct {
	sync.Map
}

func (c *TeamMap) RLoad(key TeamId) (value *Team, ok bool) {
	var v interface{}
	v, ok = c.Load(key)
	if ok {
		value = v.(*Team)
	} else {
		value = nil
	}
	return
}
