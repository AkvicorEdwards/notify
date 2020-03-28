package todo

type Tag []string

type CheckList struct {
	Msg string `json:"msg"`
	Schedule int64 `json:"schedule"`
}

type Item struct {
	Area string `json:"area"`
	Project string `json:"project"`
	DateCreate string `json:"date_create"`
	DateStart string `json:"date_start"`
	DateEnd	string `json:"date_end"`
	MsgShort string `json:"msg_short"`
	MsgFull string `json:"msg_full"`
	CheckList []CheckList `json:"check_list"`
	Tag Tag `json:"tag"`
}

// Group projects and to-dos based on different responsibilities, such as Family or Work
type Area struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
	Tag Tag `json:"tag"`
	Project Project `json:"project"`
	Item Item `json:"item"`
}

// Define a gol, then Work towards it one to-do at a time
type Project struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
	Tag Tag `json:"tag"`
	Item Item `json:"item"`
	DateCreate string `json:"date_create"`
	DateStart string `json:"date_start"`
	DateEnd	string `json:"date_end"`
}

type ToDo struct {
	Id int64 `json:"id"`
	Uid string `json:"uid"`
	Area Area `json:"area"`

}
