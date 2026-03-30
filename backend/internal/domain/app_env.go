package domain

import (
	"errors"
	"fmt"
)

type AppEnv string

const (
	AppEnvDev AppEnv = "dev"
	AppEnvPrd AppEnv = "prd"
)

var errInvalidAppEnv = errors.New("環境が不正です")

func ParseAppEnv(s string) (AppEnv, error) {
	appEnv := AppEnv(s)
	if appEnv != AppEnvDev && appEnv != AppEnvPrd {
		return "", fmt.Errorf("%w: %q", errInvalidAppEnv, s)
	}
	return appEnv, nil
}
