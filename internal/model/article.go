package model

import "errors"

type ArticleID string

var ErrInvalidArticleID = errors.New("invalid article id")

func ParseArticleID(s string) (ArticleID, error) {
	if s == "" {
		return "", ErrInvalidArticleID
	}
	return ArticleID(s), nil
}
