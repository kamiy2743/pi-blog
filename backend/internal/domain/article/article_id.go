package article

import (
	"errors"
	"fmt"
)

type ArticleID string

var errInvalidArticleID = errors.New("記事IDが不正です")

func ParseArticleID(s string) (ArticleID, error) {
	if s == "" {
		return "", fmt.Errorf("%w: %q", errInvalidArticleID, s)
	}
	return ArticleID(s), nil
}
