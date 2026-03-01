package model

import "errors"

type ArticleID string

var ErrInvalidArticleID = errors.New("記事IDが不正です")

func ParseArticleID(s string) (ArticleID, error) {
	if s == "" {
		return "", ErrInvalidArticleID
	}
	return ArticleID(s), nil
}
