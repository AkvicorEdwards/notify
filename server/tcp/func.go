package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"notify/db"
	"strconv"
	"time"
)

func cliPrintMenu() {
	fmt.Println("1. Send message to User")
	fmt.Println("2. Send message to Team")
	fmt.Println("3. Create User")
	fmt.Println("4. Create Team")
	fmt.Println("5. Add Key to User")
	fmt.Println("6. List User")
	fmt.Println("7. List Team")
	fmt.Println("8. List Connection")

}

func CliMessageSender() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("CliMessageSender recover", err)
		}
	}()
	time.Sleep(1 * time.Millisecond)

	for {
		cliPrintMenu()
		choice := GetTerminalInput()

		switch choice {
		case "1":
			cliSendMessageToUser()
		case "2":
			cliSendMessageToTeam()
		case "3":
			cliCreateUser()
		case "4":
			cliCreateTeam()
		case "5":
			cliAddKeyToUser()
		case "6":
			cliListUser()
		case "7":
			cliListTeam()
		case "8":
			cliListConnection()
		case "9":
			cliListConn()
		}
	}
}

func CreateTeam(teamId int64, teamName string, teamMemberLimit int, teamKey string, teamMember string) error {
	_, err := db.Exec("INSERT INTO user_team (team_id, team_name, team_member_limit, team_key, team_member) VALUES (?,?,?,?,?)", teamId, teamName, teamMemberLimit, teamKey, teamMember)
	return err
}

func CreateUser(uuid, uKey, nickname string, connectionLimit int, apiKey map[string]AKey) error {
	//regIPv4 := `^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`
	//regIPv6 := `^([\da-fA-F]{1,4})(:([\da-fA-F]{1,4})){7}$`
	apiK, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (uuid, u_key, nickname, connection_limit, api_key) VALUES (?,?,?,?,?)", uuid, uKey, nickname, connectionLimit, string(apiK))
	return err
}

func AddApiKey(uuid, uKey, keyId string, apiKey AKey) error {
	apiK, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE users SET api_key = JSON_INSERT(api_key, '$."`+keyId+`"', CAST(? as JSON)) WHERE uuid=? AND u_key=?`, string(apiK), uuid, uKey)
	return err
}

func SendNotification(clientId ClientId, typeCode, levelCode byte, msg []byte) {
	cli, ok := Conn.RLoad(clientId)
	if !ok {
		return
	}
	cli.DataSend <- WrapTripleCodeByte(ReqNotify, typeCode, levelCode, msg)
}

func cliSendMessageToUser() {
	check := func(str string) bool { return str == ":q" }
	var (
		err  error
		ok   bool
		user string
		con  *Connect
		id   string
		cid  int64
		tp   string
		msg  string
		cli  *Client
	)

USER:
	fmt.Println("User:")
	user = GetTerminalInput()
	if check(user) {
		return
	}
	con, ok = Connections.RLoad(user)
	if !ok {
		log.Println("User do not exist")
		goto USER
	}

ID:
	fmt.Println("ID:")
	id = GetTerminalInput()
	if check(id) {
		goto USER
	}
	cid, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("Wrong Id")
		goto ID
	}

TYPE:
	fmt.Println("1. Terminal Message")
	fmt.Println("2. Notification message")
	tp = GetTerminalInput()
	if check(tp) {
		goto ID
	}

MESSAGE:
	fmt.Println("Message:")
	msg = GetTerminalInput()
	if check(msg) {
		goto TYPE
	}

	con.Lock()
	if cli, ok = con.ConnectedClient[Int64ToClientId(cid)]; !ok {
		fmt.Println("Client Do not exist")
		goto ID
	}
	cli.RLock()
	switch tp {
	case "1":
		cli.DataSend <- WrapCodeString(ReqTerminalMessage, msg)
	case "2":
		cli.DataSend <- WrapTripleCodeString(ReqNotify, NotifyNormalMessage, NotifyNormalUrgencyLevel, msg)
	default:
		cli.DataSend <- WrapCodeString(ReqTerminalMessage, msg)
	}
	cli.RUnlock()
	con.Unlock()
	goto MESSAGE
}

func cliSendMessageToTeam() {
	check := func(str string) bool { return str == ":q" }
	var (
		err     error
		id      string
		teamId  int64
		key     string
		message string
	)

ID:
	fmt.Println("Enter Team ID:")
	id = GetTerminalInput()
	if check(id) {
		return
	}
	teamId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		goto ID
	}

KEY:
	fmt.Println("Enter Key:")
	key = GetTerminalInput()
	if check(key) {
		goto ID
	}

MESSAGE:
	fmt.Println("Enter Message:")
	message = GetTerminalInput()
	if check(message) {
		goto KEY
	}

	err = SendToTeam(Int64ToTeamId(teamId), key,
		WrapCodeString(ReqTerminalMessage, message))
	if err != nil {
		fmt.Println(err)
	}

	goto MESSAGE
}

func cliCreateUser() {
	check := func(str string) bool { return str == ":q" }
	var (
		err      error
		uuid     string
		uKey     string
		nickname string
		cLimit   string
		limit    int
	)

UUID:
	fmt.Println("Enter UUID:")
	uuid = GetTerminalInput()
	if check(uuid) {
		return
	}

KEY:
	fmt.Println("Enter User Key:")
	uKey = GetTerminalInput()
	if check(uKey) {
		goto UUID
	}

NICKNAME:
	fmt.Println("Enter Nickname")
	nickname = GetTerminalInput()
	if check(nickname) {
		goto KEY
	}

LIMIT:
	fmt.Println("Enter Connection Limit:")
	cLimit = GetTerminalInput()
	if check(cLimit) {
		goto NICKNAME
	}
	limit, err = strconv.Atoi(cLimit)
	if err != nil {
		fmt.Println("Wrong Connection Limit")
		goto LIMIT
	}

	ApiKey := make(map[string]AKey)
	err = CreateUser(uuid, uKey, nickname, limit, ApiKey)
	if err != nil {
		fmt.Println(err)
	}
}

func cliCreateTeam() {
	check := func(str string) bool { return str == ":q" }
	var (
		err             error
		teamName        string
		tLimit          string
		teamMemberLimit int
		teamKey         string
		uuid            string
		nickname        string
		teamMember      map[string]string
		tmember         []byte
	)

TEAMNAME:
	fmt.Println("Enter Team Name:")
	teamName = GetTerminalInput()
	if check(teamName) {
		return
	}

MEMBERLIMIT:
	fmt.Println("Enter Team Member Limit:")
	tLimit = GetTerminalInput()
	if check(tLimit) {
		goto TEAMNAME
	}
	teamMemberLimit, err = strconv.Atoi(tLimit)
	if err != nil {
		fmt.Println("Illegal input")
		goto MEMBERLIMIT
	}

	fmt.Println("Enter Team Key:")
	teamKey = GetTerminalInput()
	if check(teamKey) {
		goto MEMBERLIMIT
	}

	teamMember = make(map[string]string)
	fmt.Println("Add Member (:q to finished, rec to restart)")
	for {
		fmt.Println("Enter UUID")
		uuid = GetTerminalInput()
		if check(uuid) {
			break
		}
		if uuid == "rec" {
			continue
		}
		fmt.Println("Enter NickName")
		nickname = GetTerminalInput()
		if check(nickname) {
			break
		}
		if nickname == "rec" {
			continue
		}
		teamMember[uuid] = nickname
	}
	tmember, err = json.Marshal(teamMember)
	if err != nil {
		fmt.Println("Failure json Marshal")
		return
	}
	err = CreateTeam(int64(GenTeamId()), teamName, teamMemberLimit,
		teamKey, string(tmember))
	if err != nil {
		fmt.Println("Failure create team")
		return
	}
}

func cliAddKeyToUser() {
	check := func(str string) bool { return str == ":q" }

UUID:
	fmt.Println("Enter UUID:")
	uuid := GetTerminalInput()
	if check(uuid) {
		return
	}

KEY:
	fmt.Println("Enter User Key:")
	uKey := GetTerminalInput()
	if check(uKey) {
		goto UUID
	}

API:
	regIPv4 := `^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})` +
		`(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`
	regIPv6 := `^([\da-fA-F]{1,4})(:([\da-fA-F]{1,4})){7}$`
	for { // Loop 1
		perIp := PerIP{
			IP4Allow: make([]string, 0),
			IP4Deny:  make([]string, 0),
			IP6Allow: make([]string, 0),
			IP6Deny:  make([]string, 0),
		}
		fmt.Println("Enter Key ID")
		keyId := GetTerminalInput()
		if check(keyId) {
			goto KEY
		}

		fmt.Println("IP address restrictions:")
		fmt.Println("Use regular expressions")
		fmt.Println("'all' means allow all")
		fmt.Println("'rec' means re-enter")
		fmt.Println("'fin' means finished")
		fmt.Println()
		fmt.Println("Enter IPv4 Allow")
		for {
			al := GetTerminalInput()
			if check(al) {
				goto KEY
			}
			if al == "fin" {
				break
			}
			if al == "rec" {
				continue
			}
			if al == "all" {
				perIp.IP4Allow = append(perIp.IP4Allow, regIPv4)
				break
			}
			perIp.IP4Allow = append(perIp.IP4Allow, al)
			fmt.Println("Added")
		}
		fmt.Println("Enter IPv4 Deny")
		for {
			al := GetTerminalInput()
			if check(al) {
				goto KEY
			}
			if al == "fin" {
				break
			}
			if al == "rec" {
				continue
			}
			if al == "all" {
				perIp.IP4Deny = append(perIp.IP4Deny, regIPv4)
				break
			}
			perIp.IP4Deny = append(perIp.IP4Deny, al)
			fmt.Println("Added")
		}
		fmt.Println("Enter IPv6 Allow")
		for {
			al := GetTerminalInput()
			if check(al) {
				goto KEY
			}
			if al == "fin" {
				break
			}
			if al == "rec" {
				continue
			}
			if al == "all" {
				perIp.IP6Allow = append(perIp.IP6Allow, regIPv6)
				break
			}
			perIp.IP6Allow = append(perIp.IP6Allow, al)
			fmt.Println("Added")
		}
		fmt.Println("Enter IPv6 Deny")
		for {
			al := GetTerminalInput()
			if check(al) {
				goto KEY
			}
			if al == "fin" {
				break
			}
			if al == "rec" {
				continue
			}
			if al == "all" {
				perIp.IP6Deny = append(perIp.IP6Deny, regIPv6)
				break
			}
			perIp.IP6Deny = append(perIp.IP6Deny, al)
			fmt.Println("Added")
		}

	LIMIT:
		fmt.Println("Enter Connection Limit")
		climit := GetTerminalInput()
		if check(climit) {
			goto API
		}
		connectionLimit, err := strconv.Atoi(climit)
		if err != nil {
			fmt.Println(err)
			goto LIMIT
		}

	APP:
		fmt.Println("Enter App Name")
		app := GetTerminalInput()
		if check(app) {
			goto LIMIT
		}

	REMARK:
		fmt.Println("Enter Remark")
		remark := GetTerminalInput()
		if check(remark) {
			goto APP
		}

	VALID:
		fmt.Println("How long is it valid (second)")
		dLine := GetTerminalInput()
		if check(dLine) {
			goto REMARK
		}
		deadline, err := strconv.Atoi(dLine)
		if err != nil {
			fmt.Println(err)
			goto VALID
		}

		ApiKey := AKey{
			ClientID:        GenClientId(),
			ConnectionLimit: connectionLimit,
			PermittedIP:     perIp,
			App:             app,
			Remark:          remark,
			Deadline:        time.Now().Unix() + int64(deadline),
		}
		err = AddApiKey(uuid, uKey, keyId, ApiKey)
		if err != nil {
			fmt.Println(err)
		}
	} // Loop 1
}

func cliListUser() {
	fmt.Println("-------------------------------")
	dataSql := make([]SqlData, 0)
	_, err := db.Row(&dataSql, "SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make([]Data, len(dataSql))
	for k, v := range dataSql {
		data[k] = v.Transfer()
	}
	for _, v := range data {
		fmt.Printf("Nickname:[%s] UUID:[%s] KEY:[%s] ConnectionLimit:[%d]\n",
			v.Nickname, v.UUID, v.Key, v.ConnectionLimit)
		fmt.Println("API KEY:")
		for kk, vv := range v.APIKey {
			fmt.Printf("  ->Remark:[%s] APP:[%s] ID:[%v] APIKey:[%s]\n",
				vv.Remark, vv.App, vv.ClientID, kk)
			fmt.Printf("  ->ConnectionLimit:[%d] Deadline:[%d]\n",
				vv.ConnectionLimit, vv.Deadline)
			fmt.Println("  ->IPv4 Allow:")
			for _, v := range vv.PermittedIP.IP4Allow {
				fmt.Printf("      -[%s]\n", v)
			}
			fmt.Println("  ->IPv4 Deny:")
			for _, v := range vv.PermittedIP.IP4Deny {
				fmt.Printf("      -[%s]\n", v)
			}
			fmt.Println("  ->IPv6 Allow:")
			for _, v := range vv.PermittedIP.IP6Allow {
				fmt.Printf("      -[%s]\n", v)
			}
			fmt.Println("  ->IPv6 Allow:")
			for _, v := range vv.PermittedIP.IP6Deny {
				fmt.Printf("      -[%s]\n", v)
			}
		}
	}
	fmt.Println("-------------------------------")
}

func cliListTeam() {
	fmt.Println("-------------------------------")
	Teams.Range(func(key, value interface{}) bool {
		val := value.(*Team)
		fmt.Printf("NAME:[%s] ID:[%v] KEY:[%s] MemberLimit:[%d]\n",
			val.TeamName, val.Id, val.Key, val.TeamMemberLimit)
		fmt.Println("Members:")
		for k, v := range val.Member {
			fmt.Printf("  ->UUID:[%s] Nickname:[%s]\n", k, v)
		}
		fmt.Println("Connected Clients:")
		for _, v := range val.ConnectedClient {
			fmt.Printf("  ->Remark:[%s] APP:[%s] UUID:[%s] ID:[%v]\n",
				v.Remark, v.App, v.UUID, v.Id)
		}
		return true
	})
	fmt.Println("-------------------------------")
}

func cliListConnection() {
	fmt.Println("-------------------------------")
	Connections.Range(func(key, value interface{}) bool {
		val := value.(*Connect)
		fmt.Printf("Nickname:[%s] UUID:[%s] KEY:[%s] ConnectionLimit:[%d]",
			val.Nickname, val.UUID, key, val.ConnectionLimit)
		fmt.Println("Connected Client:")
		for _, v := range val.ConnectedClient {
			fmt.Printf("  ->Remark:[%s] APP:[%s] UUID:[%s] ID:[%v]\n",
				v.Remark, v.App, v.UUID, v.Id)
			for _, vv := range v.ConnectedTeam {
				fmt.Printf("    ->Team:[%s] ID:[%v]\n", vv.TeamName, vv.Id)
			}
		}
		return true
	})
	fmt.Println("-------------------------------")
}

func cliListConn() {
	fmt.Println("-------------------------------")
	Conn.Range(func(key, value interface{}) bool {
		val := value.(*Client)
		fmt.Printf("Remark:[%s] UUID:[%s] ClientID:[%v] APP:[%s]",
			val.Remark, val.UUID, key, val.App)
		return true
	})
	fmt.Println("-------------------------------")
}