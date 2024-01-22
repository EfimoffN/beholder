package types

type ConfigApp struct {
	SessionTG string
	ConfigDB  ConfigPsg
	ConfigKfk string
}

type ConfigPsg struct {
	Host     string
	User     string
	Password string
	DBname   string
	SSLmode  string
	Port     string
}

type SessionRow struct {
	SessionID   string `db:"sessionid"`
	AppID       int    `db:"appid"`
	AppHash     string `db:"apphash"`
	PhoneNumber string `db:"phonenumber"`
	Session     string `db:"sessiontxt"`
}

type LastMsgRow struct {
	MsgID      string `db:"msgid"`
	ChanneTgID int64  `db:"channelidtg"`
	LastMsgID  int    `db:"lastidmsg"`
	SessionID  string `db:"sessionid"`
}

type ChannelRow struct {
	ChannelID    string `db:"channelid"`
	ChannelTgID  int64  `db:"channelidtg"`
	ChannelName  string `db:"channelname"`
	ChannelLink  string `db:"channellink"`
	ChannelClose bool   `db:"channelclose"`
}

type ReadyRepost struct {
	ChannelTgID int64  `json:"ChannelTgID"`
	MessageLink string `json:"MessageLink"`
}
