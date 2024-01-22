package config

import (
	"errors"
	"flag"

	"github.com/EfimoffN/beholder/types"
)

var (
	ErrFlagSessionTG       = errors.New("flag session accaunt id telegram is empty")
	ErrNoAllParametersPSG  = errors.New("there are not all parameters for Postgres")
	ErrNoAllParameterKafka = errors.New("there are not parameters for Kafka")
)

func CreateConfig() (*types.ConfigApp, error) {
	cfg := types.ConfigApp{}

	flag.StringVar(&cfg.SessionTG, "a", "", "accaunt")
	flag.StringVar(&cfg.ConfigDB.User, "u", "", "userPSG")
	flag.StringVar(&cfg.ConfigDB.Password, "pa", "", "passwordPSG")
	flag.StringVar(&cfg.ConfigDB.DBname, "d", "", "dbnamePSG")
	flag.StringVar(&cfg.ConfigDB.SSLmode, "s", "", "sslmodePSG")
	flag.StringVar(&cfg.ConfigDB.Port, "po", "", "portPSG")
	flag.StringVar(&cfg.ConfigDB.Host, "h", "", "hostPSG")
	flag.StringVar(&cfg.ConfigKfk, "k", "", "kafka")

	flag.Parse()

	if cfg.SessionTG == "" {
		return nil, ErrFlagSessionTG
	}

	if cfg.ConfigDB.User == "" {
		return nil, ErrNoAllParametersPSG
	}

	if cfg.ConfigDB.Password == "" {
		return nil, ErrNoAllParametersPSG
	}

	if cfg.ConfigDB.DBname == "" {
		return nil, ErrNoAllParametersPSG
	}

	if cfg.ConfigDB.SSLmode == "" {
		return nil, ErrNoAllParametersPSG
	}

	if cfg.ConfigDB.Port == "" {
		return nil, ErrNoAllParametersPSG
	}

	if cfg.ConfigDB.Host == "" {
		return nil, ErrNoAllParametersPSG
	}

	if cfg.ConfigKfk == "" {
		return nil, ErrNoAllParameterKafka
	}

	return &cfg, nil
}
