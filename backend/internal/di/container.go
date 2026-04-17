package di

import (
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/db/ent"
	showAdminHandler "blog/internal/handler/feature/admin/show"
	searchArticleHandler "blog/internal/handler/feature/article/search"
	showArticleHandler "blog/internal/handler/feature/article/show"
	showTopHandler "blog/internal/handler/feature/top/show"
	healthHandler "blog/internal/handler/health"
	infraArticle "blog/internal/infra/article"
	infraCategory "blog/internal/infra/category"

	"github.com/romsar/gonertia/v2"
)

type Container struct {
	ShowAdminHandler     *showAdminHandler.Handler
	SearchArticleHandler *searchArticleHandler.Handler
	ShowArticleHandler   *showArticleHandler.Handler
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

	showAdminUsecase := showAdminHandler.NewUsecase(articleRepository, categoryRepository)
	searchArticleUsecase := searchArticleHandler.NewUsecase(articleRepository, categoryRepository)
	showArticleUsecase := showArticleHandler.NewUsecase(articleRepository)
	showTopUsecase := showTopHandler.NewUsecase(articleRepository, categoryRepository)

	return &Container{
		ShowAdminHandler:     showAdminHandler.NewHandler(inertiaApp, showAdminUsecase),
		SearchArticleHandler: searchArticleHandler.NewHandler(inertiaApp, searchArticleUsecase),
		ShowArticleHandler:   showArticleHandler.NewHandler(inertiaApp, showArticleUsecase),
		ShowTopHandler:       showTopHandler.NewHandler(inertiaApp, showTopUsecase),
		HealthHandler:        healthHandler.NewHandler(),
	}
}
