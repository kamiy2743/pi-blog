package article

import (
	"errors"
	"fmt"

	"blog/internal/domain"
)

type ArticleID uint32

var errInvalidArticleID = errors.New("記事IDが不正です")

func ParseArticleID(s string) (ArticleID, error) {
	id, err := domain.ParseUint32(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", errInvalidArticleID, s)
	}
	return ArticleID(id), nil
}
