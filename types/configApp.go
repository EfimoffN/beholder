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

type AcceptedPublication struct {
	ChannelTgID      int64 `json:"ChannelTgID"`
	MessageChannelID int64 `json:"MessageChannelID"`
	CreatedDate      int64 `json:"CreatedDate"`
	EditDate         int64 `json:"EditDate"`
}
