package di

import (
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/ent"
	showNotFoundHandler "blog/internal/handler/feature/notfound/show"
	showTopHandler "blog/internal/handler/feature/top/show"
	infraArticle "blog/internal/infra/article"
	infraCategory "blog/internal/infra/category"

	"github.com/romsar/gonertia/v2"
)

type Container struct {
	ShowNotFoundHandler *showNotFoundHandler.Handler
	ShowTopHandler      *showTopHandler.Handler
}

type ContainerOptions struct {
	ArticleRepository  domainArticle.ArticleRepository
	CategoryRepository domainCategory.CategoryRepository
}

func NewContainer(
	entClient *ent.Client,
	inertiaApp *gonertia.Inertia,
	options *ContainerOptions,
) *Container {
	if options == nil {
		options = &ContainerOptions{}
	}

	articleRepository := options.ArticleRepository
	if articleRepository == nil {
		articleRepository = infraArticle.NewArticleRepository(entClient)
	}

	categoryRepository := options.CategoryRepository
	if categoryRepository == nil {
		categoryRepository = infraCategory.NewCategoryRepository(entClient)
	}

	showTopUsecase := showTopHandler.NewUsecase(articleRepository, categoryRepository)

	return &Container{
		ShowNotFoundHandler: showNotFoundHandler.NewHandler(inertiaApp),
		ShowTopHandler:      showTopHandler.NewHandler(inertiaApp, showTopUsecase),
	}
}
