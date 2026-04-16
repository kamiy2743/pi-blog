package di

import (
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/ent"
	searchArticleHandler "blog/internal/handler/feature/article/search"
	showTopHandler "blog/internal/handler/feature/top/show"
	healthHandler "blog/internal/handler/health"
	infraArticle "blog/internal/infra/article"
	infraCategory "blog/internal/infra/category"

	"github.com/romsar/gonertia/v2"
)

type Container struct {
	SearchArticleHandler *searchArticleHandler.Handler
	ShowTopHandler       *showTopHandler.Handler
	HealthHandler        *healthHandler.Handler
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

	searchArticleUsecase := searchArticleHandler.NewUsecase(articleRepository, categoryRepository)
	showTopUsecase := showTopHandler.NewUsecase(articleRepository, categoryRepository)

	return &Container{
		SearchArticleHandler: searchArticleHandler.NewHandler(inertiaApp, searchArticleUsecase),
		ShowTopHandler:       showTopHandler.NewHandler(inertiaApp, showTopUsecase),
		HealthHandler:        healthHandler.NewHandler(),
	}
}
