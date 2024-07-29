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
	// Port string `required:"true"`
	// SessionTG     string `required:"true"`
	// PhoneNumber   string `required:"true"`
	// AppID         int    `required:"true"`
	// AppHASH       string `required:"true"`
	SessionOptMin int `required:"true"` // minimum value of random delay in milliseconds
	SessionOptMax int `required:"true"` // maximum value of random delay in milliseconds
	CapChan       int `default:"100"`
}

type ProducerKfk struct {
	ProducerGroup  string `required:"true"`
	ProducerTopic  string `required:"true"`
	ProducerBroker string `required:"true"`
}

// type AcceptedPublication struct {
// 	ChannelTgID int64  `json:"ChannelTgID"`
// 	MessageLink string `json:"MessageLink"`
// 	MessageID   int64  `json:"MessageID"`
// }

type AcceptedPublication2 struct {
	ChannelTgID      int64  `json:"ChannelTgID"`
	ChatTgID         int64  `json:"ChatTgID"`
	MessageChannelID int64  `json:"MessageChannelID"`
	MessageChatID    int64  `json:"MessageChatID"`
	Created          int64  `json:"Created"`
	TextMessage      string `json:"TextMessage"`
}
