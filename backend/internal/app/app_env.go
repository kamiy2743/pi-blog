package app

import (
	"fmt"
)

type AppEnv string

const (
	AppEnvDev  AppEnv = "dev"
	AppEnvStg  AppEnv = "stg"
	AppEnvPrd  AppEnv = "prd"
	AppEnvTest AppEnv = "test"
)

func ParseAppEnv(s string) (AppEnv, error) {
	appEnv := AppEnv(s)
	if appEnv != AppEnvDev && appEnv != AppEnvStg && appEnv != AppEnvPrd && appEnv != AppEnvTest {
		return "", fmt.Errorf("環境が不正です: %q", s)
	}
	return appEnv, nil
}
