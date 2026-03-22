package model

import (
	"errors"
	"fmt"
)

type AppEnv string

const (
	AppEnvDev AppEnv = "dev"
	AppEnvPrd AppEnv = "prd"
)

var ErrInvalidAppEnv = errors.New("環境が不正です")

func ParseAppEnv(s string) (AppEnv, error) {
	appEnv := AppEnv(s)
	if appEnv != AppEnvDev && appEnv != AppEnvPrd {
		return "", fmt.Errorf("%w: %q", ErrInvalidAppEnv, s)
	}
	return appEnv, nil
}
