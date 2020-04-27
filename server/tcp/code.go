package tcp

import "errors"

const MaxMessageLength int = 4 * 2000

// Request
const (
	// UUID#ApiKey
	ReqRegistration byte = 0 + iota
	//
	ReqHeartbeat
	//
	ReqTerminateTheConnection
	// TeamID#Key
	ReqEnterTeam
	// TypeCode#LevelCode#Message
	ReqNotify
	// Message
	ReqTerminalMessage
	// TeamID#Key#Req#Message
	ReqTeamMessage
	//
	ReqUserMessage
	ReqClientMessage
)

// Response
const (
	// Message
	ResRegistrationAllow byte = 50 + iota
	// ErrorCode#Message
	ResRegistrationRefuse
	ResHeartbeat
	ResTerminateTheConnection
	// Message
	ResEnterTeamAllow
	// ErrorCode#Message
	ResEnterTeamRefuse
	// Message
	ResNotifyAllow
	// ErrorCode#Message
	ResNotifyRefuse
	ResTerminalMessage
	// Message
	ResTeamMessageAllow
	// ErrorCode#Message
	ResTeamMessageRefuse
	// Message
	ResUserMessageAllow
	// ErrorCode#Message
	ResUserMessageRefuse
	// Message
	ResClientMessageAllow
	// ErrorCode#Message
	ResClientMessageRefuse
	// Code#Message
	ResError
)

// Error Code
const (
	Unknown                      byte = 100 + iota
	UnknownUser
	MaximumConnectionsExceeded
	IllegalIP
	IllegalAPIkey
	TeamDoesNotExist
	TeamMemberLimitReached
	IncorrectTeamId
	IncorrectTeamKey
	IncorrectTeamReq
	IncorrectTeamMessageReq
	IllegalRegistrationInfo
	IncorrectNotifyReq
	IllegalTeamMember
)

// Error Info
var (
	ErrorUnknown                    = errors.New("未知错误")
	ErrorUnknownUser                = errors.New("未知用户")
	ErrorMaximumConnectionsExceeded = errors.New("超出连接数量")
	ErrorIllegalIP                  = errors.New("非法IP")
	ErrorIllegalAPIkey              = errors.New("非法API密钥")
	ErrorTeamDoesNotExist           = errors.New("团队不存在")
	ErrorTeamMemberLimitReached     = errors.New("已达成员数量上限")
	ErrorIncorrectTeamId            = errors.New("错误的团队ID")
	ErrorIncorrectTeamKey           = errors.New("错误的团队密钥")
	ErrorIncorrectTeamReq           = errors.New("错误的团队请求信息")
	ErrorIncorrectTeamMessageReq    = errors.New("错误的团队信息请求")
	ErrorIllegalRegistrationInfo    = errors.New("错误的注册信息")
	ErrorIncorrectNotifyReq         = errors.New("错误的Notify请求")
	ErrorIllegalTeamMember          = errors.New("团队成员信息格式不匹配")
)

// Notify Level Code
const (
	NotifyLowUrgencyLevel      byte = iota
	NotifyNormalUrgencyLevel
	NotifyCriticalUrgencyLevel
)

// Notify Type Code
const (
	NotifyNormalMessage  byte = iota
	NotifyAndroidMessage
	NotifyToDoNotice
)


