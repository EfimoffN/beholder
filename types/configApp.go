package types

type ConfigApp struct {
	BeholderTG  SessionTG
	ConfigKfk   ProducerKfk
	ConfigDB    ConfigPsg
	SessionTGID string `required:"true"`
}

type ConfigPsg struct {
	Host     string
	User     string
	Password string
	DBname   string
	SSLmode  string
	Port     string
}

type SessionTG struct {
	SessionOptMin int `required:"true"` // minimum value of random delay in milliseconds
	SessionOptMax int `required:"true"` // maximum value of random delay in milliseconds
	CapChan       int `default:"100"`
}

type ProducerKfk struct {
	ProducerGroup  string `required:"true"`
	ProducerTopic  string `required:"true"`
	ProducerBroker string `required:"true"`
}

type AcceptedPublication2 struct {
	ChannelTgID      int64  `json:"ChannelTgID"`
	ChatTgID         int64  `json:"ChatTgID"`
	MessageChannelID int64  `json:"MessageChannelID"`
	MessageChatID    int64  `json:"MessageChatID"`
	Created          int64  `json:"Created"`
	TextMessage      string `json:"TextMessage"`
}
