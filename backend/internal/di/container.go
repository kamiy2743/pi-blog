package di

import (
	"blog/internal/db/ent"
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	destroyCategoryHandler "blog/internal/handler/feature/admin/category/destroy"
	editCategoryHandler "blog/internal/handler/feature/admin/category/edit"
	storeCategoryHandler "blog/internal/handler/feature/admin/category/store"
	updateCategoryHandler "blog/internal/handler/feature/admin/category/update"
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
	EditCategoryHandler    *editCategoryHandler.Handler
	StoreCategoryHandler   *storeCategoryHandler.Handler
	UpdateCategoryHandler  *updateCategoryHandler.Handler
	DestroyCategoryHandler *destroyCategoryHandler.Handler
	ShowAdminHandler       *showAdminHandler.Handler
	SearchArticleHandler   *searchArticleHandler.Handler
	ShowArticleHandler     *showArticleHandler.Handler
	ShowTopHandler         *showTopHandler.Handler
	HealthHandler          *healthHandler.Handler
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

	editCategoryUsecase := editCategoryHandler.NewUsecase(categoryRepository)
	showAdminUsecase := showAdminHandler.NewUsecase(articleRepository, categoryRepository)
	searchArticleUsecase := searchArticleHandler.NewUsecase(articleRepository, categoryRepository)
	showArticleUsecase := showArticleHandler.NewUsecase(articleRepository)
	showTopUsecase := showTopHandler.NewUsecase(articleRepository, categoryRepository)

	return &Container{
		EditCategoryHandler:    editCategoryHandler.NewHandler(inertiaApp, editCategoryUsecase),
		StoreCategoryHandler:   storeCategoryHandler.NewHandler(),
		UpdateCategoryHandler:  updateCategoryHandler.NewHandler(),
		DestroyCategoryHandler: destroyCategoryHandler.NewHandler(),
		ShowAdminHandler:       showAdminHandler.NewHandler(inertiaApp, showAdminUsecase),
		SearchArticleHandler:   searchArticleHandler.NewHandler(inertiaApp, searchArticleUsecase),
		ShowArticleHandler:     showArticleHandler.NewHandler(showArticleUsecase),
		ShowTopHandler:         showTopHandler.NewHandler(showTopUsecase),
		HealthHandler:          healthHandler.NewHandler(),
	}
}
