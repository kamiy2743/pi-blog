package di

import (
	"blog/internal/ent"
	showTopHandler "blog/internal/handler/top/show"
	infraArticle "blog/internal/infra/article"
	infraCategory "blog/internal/infra/category"

	"github.com/romsar/gonertia/v2"
)

type Container struct {
	ShowTopHandler *showTopHandler.Handler
}

func NewContainer(entClient *ent.Client, inertiaApp *gonertia.Inertia) *Container {
	articleRepository := infraArticle.NewArticleRepository(entClient)
	categoryRepository := infraCategory.NewCategoryRepository(entClient)
	showTopUsecase := showTopHandler.NewUsecase(articleRepository, categoryRepository)

	return &Container{
		ShowTopHandler: showTopHandler.NewHandler(inertiaApp, showTopUsecase),
	}
}
