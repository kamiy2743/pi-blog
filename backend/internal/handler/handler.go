package handler

import (
	"net/http"

	"blog/internal/app"
	"blog/internal/config"
	"blog/internal/di"
	"blog/internal/middleware"

	"blog/internal/ent"

	"github.com/romsar/gonertia/v2"
)

func NewHTTPHandler(entClient *ent.Client) (http.Handler, error) {
	inertiaApp, err := newInertiaApp()
	if err != nil {
		return nil, err
	}

	container := di.NewContainer(entClient, inertiaApp)
	mux := newMux(inertiaApp, container)

	return middleware.Chain(
		http.NewCrossOriginProtection().Handler(mux),
		middleware.NormalizePath(),
	), nil
}

func newInertiaApp() (*gonertia.Inertia, error) {
	inertiaApp, err := gonertia.NewFromFile(
		config.MustGetInertiaTemplatePath("root.html"),
		gonertia.WithSSR(config.MustGetSSRURL()),
	)
	if err != nil {
		return nil, err
	}

	inertiaApp.ShareTemplateData("faviconHref", config.MustGetTemplateFaviconHref())
	inertiaApp.ShareTemplateData("cssHref", config.MustGetTemplateCSSHref())
	inertiaApp.ShareTemplateData("useViteClient", config.MustGetAppEnv() == app.AppEnvDev)
	inertiaApp.ShareTemplateData("appScriptSrc", config.MustGetTemplateAppScriptSrc())

	return inertiaApp, nil
}
