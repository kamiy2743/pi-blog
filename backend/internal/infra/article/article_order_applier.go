package article

import (
	"fmt"

	"blog/internal/db/ent"
	entArticle "blog/internal/db/ent/article"
	domainArticle "blog/internal/domain/article"
)

func ApplyOrder(query *ent.ArticleQuery, orderBy domainArticle.OrderBy) error {
	if orderBy == "" {
		orderBy = domainArticle.OrderByDefault
	}
	switch orderBy {
	case domainArticle.OrderByLatest:
		query.Order(
			ent.Desc(entArticle.FieldUpdatedAt),
			ent.Desc(entArticle.FieldID),
		)
		return nil
	default:
		return fmt.Errorf("未対応の記事の並び順です: %s", orderBy)
	}
}
