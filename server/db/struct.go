package db

type MySQL struct {
	User string
	Password string
	Host string
	DBName string
	Charset string
	ParseTime string
	Loc string
}

type TableUser struct {
	Id			int64
	UserName	string
	UserPwd		string
	NickName	string
	Email		string
	Permissions	int64
}

func (TableUser) TableName() string {
	return "users"
}
