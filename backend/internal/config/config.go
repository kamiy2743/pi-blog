package config

import (
	"log"

	"blog/internal/domain"
)

func MustGetAppEnv() domain.AppEnv {
	appEnvRaw := mustGetEnvString("APP_ENV")
	appEnv, err := domain.ParseAppEnv(appEnvRaw)
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

func MustGetMySQLHost() string {
	return mustGetEnvString("MYSQL_HOST")
}

func MustGetMySQLPort() string {
	return mustGetEnvString("MYSQL_PORT")
}

func MustGetMySQLDatabase() string {
	return mustGetEnvString("MYSQL_DATABASE")
}

func MustGetMySQLUser() string {
	return mustGetSecretString("mysql_user")
}

func MustGetMySQLPassword() string {
	return mustGetSecretString("mysql_user_password")
}

func MustGetInertiaRootTemplatePath() string {
	return mustGetEnvString("INERTIA_ROOT_TEMPLATE_PATH")
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
