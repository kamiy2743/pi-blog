package article

import (
	"fmt"

	"blog/internal/domain"
)

type ArticleID uint32

func ParseArticleID(s string) (ArticleID, error) {
	id, err := domain.ParseUint32(s)
	if err != nil {
		return 0, fmt.Errorf("記事IDが不正です: %q", s)
	}
	return ArticleID(id), nil
}
