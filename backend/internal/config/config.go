package config

import (
	"log"

	"blog/internal/model"
)

func MustGetAppEnv() model.AppEnv {
	appEnvRaw := mustGetEnvString("APP_ENV")
	appEnv, err := model.ParseAppEnv(appEnvRaw)
	if err != nil {
		log.Fatal(err)
	}

	return appEnv
}

func MustGetPort() string {
	return mustGetEnvString("PORT")
}

func MustGetSSRURL() string {
	return mustGetEnvString("SSR_URL")
}

func MustGetInertiaRootTemplate() string {
	return mustGetEnvString("INERTIA_ROOT_TEMPLATE")
}

func MustGetTemplateFaviconHref() string {
	return mustGetEnvString("TEMPLATE_FAVICON_HREF")
}

func MustGetTemplateCSSHref() string {
	return mustGetEnvString("TEMPLATE_CSS_HREF")
}

func MustGetTemplateAppScriptSrc() string {
	return mustGetEnvString("TEMPLATE_APP_SCRIPT_SRC")
}

func MustGetAdminBasicAuthUser() string {
	return mustGetSecretString("admin_basic_auth_user")
}

func MustGetAdminBasicAuthPass() string {
	return mustGetSecretString("admin_basic_auth_pass")
}
