package tg_beholder

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/EfimoffN/beholder/types"
	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via terminal.
type termAuth struct {
	noSignUp

	phone string
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	return "", nil
}

func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")

	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(code), nil
}

func (tgc *TgBeholder) Authorize() error {
	log := zap.NewExample()

	dispatcher := tg.NewUpdateDispatcher()
	gaps := updates.New(updates.Config{
		Handler: dispatcher,
		Logger:  log.Named("gaps"),
	})

	tgc.gupMsg = gaps
	tgc.dispatcher = dispatcher

	tgOption := telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: tgc.fileStorage,
		},
		UpdateHandler: gaps,
		Logger:        log,
		// Middlewares: []telegram.Middleware{
		// 	ratelimit.New(rate.Every(time.Millisecond*200), 3),
		// },
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(gaps.Handle),
		},
	}

	client := telegram.NewClient(tgc.appID, tgc.appHASH, tgOption)

	flow := auth.NewFlow(
		termAuth{phone: tgc.phoneNumber},
		auth.SendCodeOptions{},
	)

	stop, err := bg.Connect(client)
	if err != nil {
		return errors.Wrapf(err, "can't connect")
	}

	go func() {
		for {
			if _, ok := <-tgc.done; !ok {
				_ = stop()

				return
			}
		}
	}()

	if err = client.Auth().IfNecessary(tgc.ctx, flow); err != nil {
		return errors.Wrapf(err, "failed if necessary")
	}

	tgc.client = client

	return nil
}

func (tgc *TgBeholder) Stop() {
	close(tgc.done)
}

func CreateTgBeholder(
	phoneNumber,
	appHASH,
	fileStorage string,
	appID,
	sessionOptMin,
	sessionOptMax int,
	capChan int,
	ctx context.Context) TgBeholder {

	tgClient := TgBeholder{
		phoneNumber:   phoneNumber,
		appID:         appID,
		appHASH:       appHASH,
		fileStorage:   fileStorage,
		ctx:           ctx,
		sessionOptMin: sessionOptMin, // minimum value of random delay in milliseconds
		sessionOptMax: sessionOptMax, // maximum value of random delay in milliseconds

		PostSend: make(chan types.AcceptedPublication, capChan),
	}

	return tgClient
}
