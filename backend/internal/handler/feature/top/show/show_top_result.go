package show

import (
	"blog/internal/domain/article"
	"blog/internal/domain/category"
)

type ShowTopResult struct {
	LatestArticles []article.Article
	Categories     []category.Category
}
