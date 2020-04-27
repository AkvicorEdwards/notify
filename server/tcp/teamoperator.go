package tcp

import (
	"fmt"
	"notify/db"
	"notify/encryption"
	"strconv"
	"strings"
)

// Send the request to all team members
// data format: TeamID#Key#Req#Message
// TeamID, Key, Message: URLBase64 With No Padding
func HandleSendToTeam(data []byte, cli *Client) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HandleSendToTeam recover", err)
		}
	}()
	// TeamID  Key  Req#Message
	info := strings.SplitN(string(data), "#", 2)
	if len(info) < 3 {
		cli.DataSend <- WrapDoubleCodeString(ResTeamMessageRefuse, IncorrectTeamMessageReq,
			ErrorIncorrectTeamMessageReq.Error())
		return
	}
	info[0] = encryption.Base64StringToString(info[0])
	info[1] = encryption.Base64StringToString(info[1])
	id, err := strconv.ParseInt(info[0], 10, 64)
	if err != nil {
		cli.DataSend <- WrapDoubleCodeString(ResTeamMessageRefuse, IncorrectTeamId,
			ErrorIncorrectTeamId.Error())
		return
	}
	err = SendToTeam(Int64ToTeamId(id), info[1], []byte(info[2]))
	switch err {
	case ErrorTeamDoesNotExist:
		cli.DataSend <- WrapDoubleCodeString(ResTeamMessageRefuse, TeamDoesNotExist,
			ErrorTeamDoesNotExist.Error())
	case ErrorIncorrectTeamKey:
		cli.DataSend <- WrapDoubleCodeString(ResTeamMessageRefuse, IncorrectTeamKey,
			ErrorIncorrectTeamKey.Error())
	case nil:
		cli.DataSend <- WrapCodeString(ResTeamMessageAllow, "团队信息已发送")
	default:
		cli.DataSend <- WrapDoubleCodeString(ResTeamMessageRefuse, Unknown,
			ErrorUnknown.Error())
	}
}

func SendToTeam(id TeamId, key string, data []byte) error {
	team, ok := Teams.RLoad(id)
	if !ok {
		return ErrorTeamDoesNotExist
	}
	if key != team.Key {
		return ErrorIncorrectTeamKey
	}
	team.Lock()
	for _, v := range team.ConnectedClient {
		v.DataSend <- data
	}
	team.Unlock()
	return nil
}

// data format: TeamID#Key
// TeamID, Key: URLBase64 With No Padding
func HandleEnterTeam(data []byte, cli *Client) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HandleEnterTeam recover", err)
		}
	}()

	info := strings.Split(string(data), "#")

	if len(info) < 2 {
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, IncorrectTeamReq,
			ErrorIncorrectTeamReq.Error())
		return
	}
	info[0] = encryption.Base64StringToString(info[0])
	info[1] = encryption.Base64StringToString(info[1])

	id, err := strconv.ParseInt(info[0], 10, 64)
	if err != nil {
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, IncorrectTeamId,
			ErrorIncorrectTeamId.Error())
		return
	}
	err = EnterTeam(Int64ToTeamId(id), info[1], cli)
	switch err {
	case ErrorTeamDoesNotExist:
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, TeamDoesNotExist,
			ErrorTeamDoesNotExist.Error())
	case ErrorIllegalTeamMember:
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, IllegalTeamMember,
			ErrorIllegalTeamMember.Error())
	case ErrorTeamMemberLimitReached:
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, TeamMemberLimitReached,
			ErrorTeamMemberLimitReached.Error())
	case ErrorIncorrectTeamKey:
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, IncorrectTeamKey,
			ErrorIncorrectTeamKey.Error())
	case nil:
		cli.DataSend <- WrapCodeString(ResEnterTeamAllow, "成功进入团队")
	default:
		cli.DataSend <- WrapDoubleCodeString(ResEnterTeamRefuse, Unknown,
			ErrorUnknown.Error())
	}
}

func EnterTeam(id TeamId, key string, cli *Client) error {
	team, ok := Teams.RLoad(id)
	if !ok {
		return ErrorTeamDoesNotExist
	}
	if _, ok := team.Member[cli.UUID]; !ok {
		return ErrorIllegalTeamMember
	}
	if len(team.ConnectedClient) >= team.TeamMemberLimit {
		return ErrorTeamMemberLimitReached
	}
	if key != team.Key {
		return ErrorIncorrectTeamKey
	}
	team.Lock()
	team.ConnectedClient[cli.Id] = cli
	team.Unlock()
	cli.Lock()
	cli.ConnectedTeam[id] = team
	cli.Unlock()
	return nil
}

func TeamLoader() {
	All := make([]SqlTeam, 0)
	_, err := db.Row(&All, "SELECT * FROM user_team LIMIT 100")
	if err != nil {
		fmt.Println("Team Loader Error in MySQL")
		return
	}
	for _, v := range All {
		va := v.Transfer()
		Teams.Store(va.Id, &va)
	}
}
