package tcp

import (
	"fmt"
	"net"
	"notify/db"
	"regexp"
	"time"
)

func AddConnections(uuid string, ip string, apiKey string, con net.Conn) (*Client, bool) {
	// 检测内存中是否已经存在用户数据
	// 若不存在则从数据库中查询用户数据并放入内存
	c, ok := Connections.RLoad(uuid)
	if !ok {
		data := &SqlData{}
		_, err := db.Row(data, "SELECT uuid, nickname, connection_limit " +
			" FROM users WHERE uuid=? LIMIT 1", uuid)
		if err != nil {
			fmt.Println(uuid)
			_, _ = con.Write(WrapDoubleCodeString(ResRegistrationRefuse, UnknownUser,
				ErrorUnknownUser.Error()))
			return nil, false
		}
		c = NewConnect(data.UUID, data.Nickname, data.ConnectionLimit)
		Connections.Store(uuid, c)
	}
	id, err := c.AddClient(uuid, ip, apiKey, con)
	if err != nil {
		switch err {
		case ErrorMaximumConnectionsExceeded:
			//fmt.Println("Maximum Connections Exceeded AddConnections")
			_, _ = con.Write(WrapDoubleCodeString(ResRegistrationRefuse,
				MaximumConnectionsExceeded, ErrorMaximumConnectionsExceeded.Error()))
		case ErrorIllegalIP:
			//fmt.Println("Illegal IP in AddConnections")
			_, _ = con.Write(WrapDoubleCodeString(ResRegistrationRefuse,
				IllegalIP, ErrorIllegalIP.Error()))
		case ErrorIllegalAPIkey:
			//fmt.Println("Illegal API Key in AddConnections")
			_, _ = con.Write(WrapDoubleCodeString(ResRegistrationRefuse,
				IllegalAPIkey, ErrorIllegalAPIkey.Error()))
		default:
			//fmt.Println("Error in AddConnections")
			_, _ = con.Write(WrapDoubleCodeString(ResRegistrationRefuse,
				Unknown, ErrorUnknown.Error()))
		}
		return nil, false
	}



	//fmt.Println(uuid, "Successful in AddConnections")
	_, _ = con.Write(WrapCodeString(ResRegistrationAllow, "注册成功"))

	// Ensure the return value is valid
	// 'return' will be executed before 'defer'
	c.RLock()
	defer c.RUnlock()
	if x, ok := c.ConnectedClient[id]; ok {
		Conn.Store(id, x)
		return x, true
	}
	return nil, false
}

func (c *Connect) AddClient(uuid, ip, apiKey string, con net.Conn) (ClientId, error) {
	var data Data
	var k AKey
	var id ClientId
	ok := false

	// Check APIKey
	if !func() bool {
		func() {
			//d := &SqlData{}
			var d SqlData
			_, err := db.Row(&d, "SELECT api_key FROM users WHERE uuid=? LIMIT 1", uuid)
			if err != nil {
				fmt.Println(err)
				return
			}
			data = d.Transfer()
		}()
		if k, ok = data.APIKey[apiKey]; ok {
			if time.Now().Unix() <= k.Deadline {
				id = k.ClientID
				return true
			}
		}
		return false
	}() {
		return id, ErrorIllegalAPIkey
	}


	// Check the number of connections
	if len(c.ConnectedClient) >= c.ConnectionLimit ||
		len(c.ConnectedClient) >= k.ConnectionLimit {
		return id, ErrorMaximumConnectionsExceeded
	}

	// Check IP
	if !func() bool {
		regIPv4 := `^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})` +
			`(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`
		regIPv6 := `^([\da-fA-F]{1,4})(:([\da-fA-F]{1,4})){7}$`
		if matched, err := regexp.MatchString(regIPv4, ip); err == nil && matched {
			for _, v := range k.PermittedIP.IP4Deny {
				if matched, err := regexp.MatchString(v, ip); err == nil && matched {
					return false
				}
			}
			for _, v := range k.PermittedIP.IP4Allow {
				if matched, err := regexp.MatchString(v, ip); err == nil && matched {
					return true
				}
			}
		} else if matched, err := regexp.MatchString(regIPv6, ip); err == nil && matched {
			for _, v := range k.PermittedIP.IP6Deny {
				if matched, err := regexp.MatchString(v, ip); err == nil && matched {
					return false
				}
			}
			for _, v := range k.PermittedIP.IP6Allow {
				if matched, err := regexp.MatchString(v, ip); err == nil && matched {
					return true
				}
			}
		}
		return false
	}() {
		return id, ErrorIllegalIP
	}

	// End of inspection, no problems found

	if k, ok := c.ConnectedClient[id]; ok {
		k.DataSend <- WrapCode(ReqTerminateTheConnection)
		time.Sleep(1 * time.Second)
		k.Termination <- true
	}

	time.Sleep(1 * time.Second)

	for {
		if _, ok := c.ConnectedClient[id]; ok {
			time.Sleep(1 * time.Second)
		} else {
			c.Lock()
			c.ConnectedClient[id] = NewClient(id, uuid, k.App, k.Remark, k.Deadline, con)
			c.Unlock()
			break
		}
	}

	return id, nil
}
