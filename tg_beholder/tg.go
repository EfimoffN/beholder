package tg_beholder

import (
	"context"

	"github.com/EfimoffN/beholder/types"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/rs/zerolog"
)

type TgBeholder struct {
	Logger        zerolog.Logger
	client        *telegram.Client
	done          chan (struct{})
	phoneNumber   string
	appID         int
	appHASH       string
	fileStorage   string
	ctx           context.Context
	sessionOptMin int // minimum value of random delay in milliseconds
	sessionOptMax int // maximum value of random delay in milliseconds

	gupMsg     *updates.Manager
	dispatcher tg.UpdateDispatcher

	PostSend chan (types.AcceptedPublication)
}

type SessionData struct {
	Session     SessionTG
	AccountData AccountData
}
type AccountData struct {
	AppID             int
	AppHash           string
	Phone             string
	Password          string
	Session           string
	Status            string
	StatusDescription string
}

type SessionTG struct {
	Version int `json:"Version"`
	Data    struct {
		Config    ConfigTG `json:"Config"`
		Dc        int      `json:"DC"`
		Addr      string   `json:"Addr"`
		AuthKey   string   `json:"AuthKey"`
		AuthKeyID string   `json:"AuthKeyID"`
		Salt      int64    `json:"Salt"`
	} `json:"Data"`
}

type ConfigTG struct {
	BlockedMode     bool        `json:"BlockedMode"`
	PFSEnabled      bool        `json:"PFSEnabled"`
	ForceTryIpv6    bool        `json:"ForceTryIpv6"`
	Date            int         `json:"Date"`
	Expires         int         `json:"Expires"`
	TestMode        bool        `json:"TestMode"`
	ThisDC          int         `json:"ThisDC"`
	DCOptions       []DCOptions `json:"DCOptions"`
	DCTxtDomainName string      `json:"DCTxtDomainName"`
	TmpSessions     int         `json:"TmpSessions"`
	WebfileDCID     int         `json:"WebfileDCID"`
}

type DCOptions struct {
	Flags             int         `json:"Flags"`
	Ipv6              bool        `json:"Ipv6"`
	MediaOnly         bool        `json:"MediaOnly"`
	TCPObfuscatedOnly bool        `json:"TCPObfuscatedOnly"`
	Cdn               bool        `json:"CDN"`
	Static            bool        `json:"Static"`
	ThisPortOnly      bool        `json:"ThisPortOnly"`
	ID                int         `json:"ID"`
	IPAddress         string      `json:"IPAddress"`
	Port              int         `json:"Port"`
	Secret            interface{} `json:"Secret"`
}
