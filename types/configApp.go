package types

type ConfigApp struct {
	BeholderTG SessionTG
	ConfigKfk  ProducerKfk
}

type SessionTG struct {
	Port          string `required:"true"`
	SessionTG     string `required:"true"`
	PhoneNumber   string `required:"true"`
	AppID         int    `required:"true"`
	AppHASH       string `required:"true"`
	SessionOptMin int    `required:"true"` // minimum value of random delay in milliseconds
	SessionOptMax int    `required:"true"` // maximum value of random delay in milliseconds
	CapChan       int    `default:"100"`
}

type ProducerKfk struct {
	ProducerGroup  string `required:"true"`
	ProducerTopic  string `required:"true"`
	ProducerBroker string `required:"true"`
}

type AcceptedPublication struct {
	ChannelTgID int64  `json:"ChannelTgID"`
	MessageLink string `json:"MessageLink"`
	MessageID   int64  `json:"MessageID"`
}
