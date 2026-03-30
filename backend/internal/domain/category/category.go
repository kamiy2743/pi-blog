package category

import (
	"errors"
	"fmt"
	"strings"
)

type Category struct {
	ID   CategoryID
	Name string
}

type CreateCategoryInput struct {
	Name string
}

var (
	errInvalidCategory   = errors.New("カテゴリが不正です")
	errEmptyCategoryName = errors.New("カテゴリ名は必須です")
)

func (c Category) Validate() error {
	if err := validateContent(c.Name); err != nil {
		return err
	}
	return nil
}

func (c CreateCategoryInput) Validate() error {
	if err := validateContent(c.Name); err != nil {
		return err
	}
	return nil
}

func validateContent(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: %w", errInvalidCategory, errEmptyCategoryName)
	}
	return nil
}
