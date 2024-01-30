package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/EfimoffN/beholder/types"
)

var (
	ErrEnvPortGRPCEmpty    = errors.New("env grpc port is empty")
	ErrEnvSessionFileEmpty = errors.New("env session file is empty")
	ErrEnvPhoneNumberEmpty = errors.New("env phone number empty")
	ErrEnvAppIDEmpty       = errors.New("env app id is empty")
	ErrEnvAppHASHEmpty     = errors.New("env app hash empty")
	ErrEnvSessOptMinEmpty  = errors.New("env session opt min is empty")
	ErrEnvSessOptMaxEmpty  = errors.New("env session opt max is empty")
	ErrEnvBrokerProducer   = errors.New("env broker producer opt is empty")
	ErrEnvTopicProducer    = errors.New("env topic producer opt is empty")
	ErrEnvGroupProducer    = errors.New("env group producer opt is empty")
)

func CreateConfig() (*types.ConfigApp, error) {
	cfg := types.ConfigApp{}

	beholderTG, err := getBeholderTGConfig()
	if err != nil {
		return nil, err
	}

	cfg.BeholderTG = *beholderTG

	configKfk, err := getProducerKfk()
	if err != nil {
		return nil, err
	}

	cfg.ConfigKfk = *configKfk

	return &cfg, nil
}

func getBeholderTGConfig() (*types.SessionTG, error) {
	beholderTG := types.SessionTG{}

	beholderTG.Port = os.Getenv("PORT")
	if beholderTG.Port == "" {
		return nil, ErrEnvPortGRPCEmpty
	}

	beholderTG.SessionTG = os.Getenv("SESSION_TG")
	if beholderTG.SessionTG == "" {
		return nil, ErrEnvSessionFileEmpty
	}

	beholderTG.PhoneNumber = os.Getenv("PHONE_NUMBER")
	if beholderTG.SessionTG == "" {
		return nil, ErrEnvPhoneNumberEmpty
	}

	appID := os.Getenv("APP_ID")
	if appID == "" {
		return nil, ErrEnvAppIDEmpty
	}

	ad, err := strconv.Atoi(appID)
	if err != nil {
		return nil, err
	}

	beholderTG.AppID = ad

	beholderTG.AppHASH = os.Getenv("APP_HASH")
	if beholderTG.AppHASH == "" {
		return nil, ErrEnvAppHASHEmpty
	}

	sessionOptMax := os.Getenv("SESS_OPT_MIN")
	if sessionOptMax == "" {
		return nil, ErrEnvSessOptMaxEmpty
	}

	sMax, err := strconv.Atoi(sessionOptMax)
	if err != nil {
		return nil, err
	}

	beholderTG.SessionOptMax = sMax

	sessionOptMin := os.Getenv("SESS_OPT_MAX")
	if sessionOptMin == "" {
		return nil, ErrEnvSessOptMinEmpty
	}

	sMin, err := strconv.Atoi(sessionOptMin)
	if err != nil {
		return nil, err
	}

	beholderTG.SessionOptMin = sMin

	return &beholderTG, nil
}

func getProducerKfk() (*types.ProducerKfk, error) {
	producerKfk := types.ProducerKfk{}

	if producerKfk.ProducerBroker == os.Getenv("BROKER") {
		return nil, ErrEnvBrokerProducer
	}

	if producerKfk.ProducerBroker == os.Getenv("TOPIC") {
		return nil, ErrEnvTopicProducer
	}

	if producerKfk.ProducerBroker == os.Getenv("GROUP") {
		return nil, ErrEnvGroupProducer
	}

	return &producerKfk, nil
}
