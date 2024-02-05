package types

type PrjSessionRow struct {
	SessionID   string `db:"sessionid"`
	Appid       string `db:"appid"`
	AppHash     string `db:"apphash"`
	PhoneNumber string `db:"phonenumber"`
	Sessiontxt  string `db:"sessiontxt"`
}
