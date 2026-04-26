package app

import (
	"fmt"
)

type AppEnv string

const (
	AppEnvDev  AppEnv = "dev"
	AppEnvPrd  AppEnv = "prd"
	AppEnvTest AppEnv = "test"
)

func ParseAppEnv(s string) (AppEnv, error) {
	appEnv := AppEnv(s)
	if appEnv != AppEnvDev && appEnv != AppEnvPrd && appEnv != AppEnvTest {
		return "", fmt.Errorf("環境が不正です: %q", s)
	}
	return appEnv, nil
}
