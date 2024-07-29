package types

type PrjClient struct {
	ClientID   string `db:"clientid"`
	ClientName string `db:"clientname"`
	ClientDesc string `db:"clientdesc"`
	ClientOff  bool   `db:"clientoff"`
}

type ChannelRow struct {
	ChannelID    string `db:"channelid"`
	ChannelTgID  int64  `db:"channelidtg"`
	ChannelName  string `db:"channelname"`
	ChannelLink  string `db:"channellink"`
	ChannelClose bool   `db:"channelclose"`
}

type PrjSessionRow struct {
	SessionID   string `db:"sessionid"`
	AppID       int    `db:"appid"`
	AppHash     string `db:"apphash"`
	PhoneNumber string `db:"phonenumber"`
	Sessiontxt  string `db:"sessiontxt"`
}

type RefClientChannel struct {
	RefID          string `db:"refid"`
	ClientID       int64  `db:"clientid"`
	ChannelID      int64  `db:"channelid"`
	ExpirationDate string `db:"expirationdate"`
}

type RefClientSession struct {
	RefID          string `db:"refid"`
	ClientID       int64  `db:"clientid"`
	SessionID      int64  `db:"sessionid"`
	ExpirationDate string `db:"expirationdate"`
}
