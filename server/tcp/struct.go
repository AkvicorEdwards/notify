package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type MSG struct {
	App   string `json:"app"`
	Title string `json:"title"`
	Msg   string `json:"msg"`
}

type SqlTeam struct {
	Id              int    `gorm:"column:id"`
	TeamId          int64  `gorm:"column:team_id"`
	TeamName        string `gorm:"column:team_name"`
	TeamMemberLimit int    `gorm:"column:team_member_limit"`
	Key             string `gorm:"column:team_key"`
	Member          string `gorm:"column:team_member"`
}

func (d *SqlTeam) Transfer() (a Team) {
	a.Id = Int64ToTeamId(d.TeamId)
	a.TeamName = d.TeamName
	a.TeamMemberLimit = d.TeamMemberLimit
	a.Key = d.Key
	err := json.Unmarshal([]byte(d.Member), &a.Member)
	if err != nil {
		fmt.Println(err)
	}
	a.ConnectedClient = make(map[ClientId]*Client)
	a.RWMutex = sync.RWMutex{}
	return
}
type Team struct {
	Id              TeamId
	TeamName        string
	TeamMemberLimit int
	Key             string
	// UUID => Nickname
	Member          map[string]string
	// ClientID => *Client
	ConnectedClient map[ClientId]*Client
	sync.RWMutex
}

func NewTeam(id TeamId, name string, limit int, key string) *Team {
	return &Team{
		Id:              id,
		TeamName:        name,
		TeamMemberLimit: limit,
		Key:             key,
		Member:          make(map[string]string),
		ConnectedClient: make(map[ClientId]*Client),
		RWMutex:         sync.RWMutex{},
	}
}

type SqlData struct {
	UUID            string `gorm:"column:uuid"`
	Key             string `gorm:"column:u_key"`
	Nickname        string `gorm:"column:nickname"`
	ConnectionLimit int    `gorm:"column:connection_limit"`
	APIKey          string `gorm:"column:api_key"`
}

func (d *SqlData) Transfer() (a Data) {
	a.UUID = d.UUID
	a.Key = d.Key
	a.Nickname = d.Nickname
	a.ConnectionLimit = d.ConnectionLimit
	err := json.Unmarshal([]byte(d.APIKey), &a.APIKey)
	if err != nil {
		fmt.Println(err)
	}
	return
}

type Data struct {
	UUID            string          `json:"uuid"`
	Key             string          `json:"key"`
	Nickname        string          `json:"nickname"`
	ConnectionLimit int             `json:"connection_limit"`
	// API KEY => AKey
	APIKey          map[string]AKey `json:"api_key"`
}

type AKey struct {
	ClientID        ClientId `json:"client_id"`
	ConnectionLimit int    `json:"connection_limit"`
	PermittedIP     PerIP  `json:"permitted_ip"`
	App             string `json:"app"`
	Remark          string `json:"remark"`
	Deadline        int64  `json:"deadline"`
}

type PerIP struct {
	IP4Allow []string `json:"ipv4_allow"`
	IP4Deny  []string `json:"ipv4_deny"`
	IP6Allow []string `json:"ipv6_allow"`
	IP6Deny  []string `json:"ipv6_deny"`
}

type Client struct {
	Id                    ClientId
	UUID                  string
	App                   string
	Remark                string
	Deadline              int64
	Connection            net.Conn
	DataReceived          chan []byte
	DataSend              chan []byte
	HeartbeatDataReceived chan bool
	HeartbeatDataSent     chan bool
	Termination           chan bool
	Worker                map[string]chan bool
	ConnectedTeam         map[TeamId]*Team
	sync.RWMutex
}

func NewClient(id ClientId, uuid, app, remark string, deadline int64, con net.Conn) *Client {
	return &Client{
		Id:                    id,
		UUID:                  uuid,
		App:                   app,
		Remark:                remark,
		Deadline:              deadline,
		Connection:            con,
		DataReceived:          make(chan []byte, 2),
		DataSend:              make(chan []byte, 2),
		HeartbeatDataReceived: make(chan bool, 1),
		HeartbeatDataSent:     make(chan bool, 1),
		Termination:           make(chan bool),
		Worker:                make(map[string]chan bool),
		ConnectedTeam:         make(map[TeamId]*Team, 0),
		RWMutex:               sync.RWMutex{},
	}
}

type Connect struct {
	UUID            string
	Nickname        string
	ConnectionLimit int
	ConnectedClient map[ClientId]*Client
	sync.RWMutex
}

func NewConnect(uuid, nickname string, limit int) *Connect {
	return &Connect{
		UUID:            uuid,
		Nickname:        nickname,
		ConnectionLimit: limit,
		ConnectedClient: make(map[ClientId]*Client),
		RWMutex:         sync.RWMutex{},
	}
}

