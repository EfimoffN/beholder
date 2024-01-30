package types

type ConfigApp struct {
	BeholderTG SessionTG
	ConfigKfk  ProducerKfk
}

type ReadyRepost struct {
	ChannelTgID int64  `json:"ChannelTgID"`
	MessageLink string `json:"MessageLink"`
}

type SessionTG struct {
	Port          string
	SessionTG     string
	PhoneNumber   string
	AppID         int
	AppHASH       string
	SessionOptMin int // minimum value of random delay in milliseconds
	SessionOptMax int // maximum value of random delay in milliseconds
}

type ProducerKfk struct {
	ProducerGroup  string
	ProducerTopic  string
	ProducerBroker string
}
