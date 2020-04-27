package tcp

import "errors"

const MaxMessageLength int = 4 * 2000

// Request
const (
	// UUID#ApiKey
	ReqRegistration byte = 0
	//
	ReqHeartbeat byte = 1
	//
	ReqTerminateTheConnection byte = 2
	// TeamID#Key
	ReqEnterTeam byte = 3
	// TypeCode#LevelCode#Message
	ReqNotify byte = 4
	// Message
	ReqTerminalMessage byte = 5
	// TeamID#Key#Req#Message
	ReqTeamMessage byte = 6
)

// Response
const (
	// Message
	ResRegistrationAllow byte = 50
	// ErrorCode#Message
	ResRegistrationRefuse     byte = 51
	ResHeartbeat              byte = 52
	ResTerminateTheConnection byte = 53
	// Message
	ResEnterTeamAllow byte = 54
	// ErrorCode#Message
	ResEnterTeamRefuse byte = 54
	ResNotify          byte = 55
	ResTerminalMessage byte = 56
	// Message
	ResTeamMessageAllow byte = 57
	// ErrorCode#Message
	ResTeamMessageRefuse byte = 57
	// Code#Message
	ResError byte = 58
)

// Notify Level Code
const (
	NotifyLowUrgencyLevel      byte = 0
	NotifyNormalUrgencyLevel   byte = 1
	NotifyCriticalUrgencyLevel byte = 2
)

// Notify Type Code
const (
	NotifyNormalMessage  byte = 0
	NotifyAndroidMessage byte = 1
	NotifyToDoNotice     byte = 2
)

// Error
const (
	Error                      byte = 100
	IllegalRegistrationInfo    byte = 101
	MaximumConnectionsExceeded byte = 102
	IllegalIP                  byte = 103
	IllegalAPIkey              byte = 104
	TeamDoesNotExist           byte = 105
	TeamMemberLimitReached     byte = 106
	IncorrectTeamId            byte = 107
	IncorrectTeamKey           byte = 108
	IncorrectTeamReq           byte = 109
	IncorrectTeamMessageReq    byte = 110
	UnknownUser                byte = 111
)

var (
	ErrorUnknown                    = errors.New("未知错误")
	ErrorUnknownUser                = errors.New("未知用户")
	ErrorMaximumConnectionsExceeded = errors.New("超出连接数量")
	ErrorIllegalIP                  = errors.New("非法IP")
	ErrorIllegalAPIkey              = errors.New("非法API密钥")
	ErrorIllegalTeamMember          = errors.New("团队成员信息格式不匹配")
	ErrorTeamDoesNotExist           = errors.New("团队不存在")
	ErrorTeamMemberLimitReached     = errors.New("已达成员数量上限")
	ErrorIncorrectTeamId            = errors.New("错误的团队ID")
	ErrorIncorrectTeamKey           = errors.New("错误的团队密钥")
	ErrorIncorrectTeamReq           = errors.New("错误的团队请求信息")
	ErrorIncorrectTeamMessageReq    = errors.New("错误的团队信息请求")
	ErrorIllegalRegistrationInfo    = errors.New("错误的注册信息")
)
