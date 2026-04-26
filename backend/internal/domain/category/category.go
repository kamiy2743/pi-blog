package category

import (
	"errors"
	"strings"
)

type Category struct {
	ID   CategoryID
	Name string
}

type CreateCategoryInput struct {
	Name string
}

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
		return errors.New("カテゴリ名は必須です")
	}
	return nil
}
