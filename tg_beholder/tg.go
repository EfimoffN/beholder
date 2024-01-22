package tgcrawl

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/rs/zerolog"
)

type TgCrawler struct {
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

func InitTgCrawler(
	logger zerolog.Logger,
	phone string,
	fileSt string,
	appID int,
	appHASH string,
	ctx context.Context,
	sessionOptMin,
	sessionOptMax int,
) *TgCrawler {
	crawler := &TgCrawler{
		Logger:        logger,
		done:          make(chan struct{}),
		phoneNumber:   phone,
		fileStorage:   fileSt,
		ctx:           ctx,
		appID:         appID,
		appHASH:       appHASH,
		sessionOptMin: sessionOptMin,
		sessionOptMax: sessionOptMax,
	}

	return crawler
}
