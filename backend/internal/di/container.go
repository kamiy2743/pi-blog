package di

import (
	"blog/internal/db/ent"
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	createArticleHandler "blog/internal/handler/feature/admin/article/create"
	destroyArticleHandler "blog/internal/handler/feature/admin/article/destroy"
	editArticleHandler "blog/internal/handler/feature/admin/article/edit"
	storeArticleHandler "blog/internal/handler/feature/admin/article/store"
	updateArticleHandler "blog/internal/handler/feature/admin/article/update"
	destroyCategoryHandler "blog/internal/handler/feature/admin/category/destroy"
	editCategoryHandler "blog/internal/handler/feature/admin/category/edit"
	storeCategoryHandler "blog/internal/handler/feature/admin/category/store"
	updateCategoryHandler "blog/internal/handler/feature/admin/category/update"
	showAdminHandler "blog/internal/handler/feature/admin/show"
	searchArticleHandler "blog/internal/handler/feature/article/search"
	showArticleHandler "blog/internal/handler/feature/article/show"
	showTopHandler "blog/internal/handler/feature/top/show"
	infraArticle "blog/internal/infra/article"
	infraCategory "blog/internal/infra/category"
)

type Container struct {
	CreateArticleHandler   *createArticleHandler.Handler
	DestroyArticleHandler  *destroyArticleHandler.Handler
	EditArticleHandler     *editArticleHandler.Handler
	StoreArticleHandler    *storeArticleHandler.Handler
	UpdateArticleHandler   *updateArticleHandler.Handler
	DestroyCategoryHandler *destroyCategoryHandler.Handler
	EditCategoryHandler    *editCategoryHandler.Handler
	StoreCategoryHandler   *storeCategoryHandler.Handler
	UpdateCategoryHandler  *updateCategoryHandler.Handler
	ShowAdminHandler       *showAdminHandler.Handler
	SearchArticleHandler   *searchArticleHandler.Handler
	ShowArticleHandler     *showArticleHandler.Handler
	ShowTopHandler         *showTopHandler.Handler
}

type ContainerOptions struct {
	ArticleRepository  domainArticle.ArticleRepository
	CategoryRepository domainCategory.CategoryRepository
}

func NewContainer(entClient *ent.Client, options *ContainerOptions) *Container {
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
	storeCategoryUsecase := storeCategoryHandler.NewUsecase(categoryRepository)
	updateCategoryUsecase := updateCategoryHandler.NewUsecase(categoryRepository)
	destroyCategoryUsecase := destroyCategoryHandler.NewUsecase(categoryRepository)
	createArticleUsecase := createArticleHandler.NewUsecase(categoryRepository)
	destroyArticleUsecase := destroyArticleHandler.NewUsecase(articleRepository)
	editArticleUsecase := editArticleHandler.NewUsecase(articleRepository, categoryRepository)
	storeArticleUsecase := storeArticleHandler.NewUsecase(articleRepository, categoryRepository)
	updateArticleUsecase := updateArticleHandler.NewUsecase(articleRepository, categoryRepository)
	showAdminUsecase := showAdminHandler.NewUsecase(articleRepository, categoryRepository)
	searchArticleUsecase := searchArticleHandler.NewUsecase(articleRepository, categoryRepository)
	showArticleUsecase := showArticleHandler.NewUsecase(articleRepository)
	showTopUsecase := showTopHandler.NewUsecase(articleRepository, categoryRepository)

	return &Container{
		CreateArticleHandler:   createArticleHandler.NewHandler(createArticleUsecase),
		DestroyArticleHandler:  destroyArticleHandler.NewHandler(destroyArticleUsecase),
		EditArticleHandler:     editArticleHandler.NewHandler(editArticleUsecase),
		StoreArticleHandler:    storeArticleHandler.NewHandler(storeArticleUsecase),
		UpdateArticleHandler:   updateArticleHandler.NewHandler(updateArticleUsecase),
		DestroyCategoryHandler: destroyCategoryHandler.NewHandler(destroyCategoryUsecase),
		EditCategoryHandler:    editCategoryHandler.NewHandler(editCategoryUsecase),
		StoreCategoryHandler:   storeCategoryHandler.NewHandler(storeCategoryUsecase),
		UpdateCategoryHandler:  updateCategoryHandler.NewHandler(updateCategoryUsecase),
		ShowAdminHandler:       showAdminHandler.NewHandler(showAdminUsecase),
		SearchArticleHandler:   searchArticleHandler.NewHandler(searchArticleUsecase),
		ShowArticleHandler:     showArticleHandler.NewHandler(showArticleUsecase),
		ShowTopHandler:         showTopHandler.NewHandler(showTopUsecase),
	}
}
