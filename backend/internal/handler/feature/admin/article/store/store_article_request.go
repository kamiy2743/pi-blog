package store

import (
	"net/http"
	"strings"
	"time"

	"blog/internal/domain"
	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
	"blog/internal/handler/handlerrequest"
	"blog/internal/handler/validator"
)

type request struct {
	Title          string   `json:"title" validate:"required,notblank,max=255"`
	Body           string   `json:"body" validate:"required,notblank"`
	IsPublished    string   `json:"isPublished" validate:"required,bool"`
	PublishStartAt string   `json:"publishStartAt" validate:"omitempty,datetime,datetime_lt=PublishEndAt"`
	PublishEndAt   string   `json:"publishEndAt" validate:"omitempty,datetime"`
	CategoryIDs    []string `json:"categoryIds"`
}

type input struct {
	Title          string
	Body           string
	IsPublished    bool
	PublishStartAt *time.Time
	PublishEndAt   *time.Time
	CategoryIDs    []category.CategoryID
}

func toInput(r *http.Request) (input, *handlererror.ValidationError) {
	req := request{}
	handlerrequest.DecodeJSONForm(r, &req)
	validationError := validator.Validate(req, getValidationMessage)

	categoryIDs, parseCategoryIDsError := parseCategoryIDs(req.CategoryIDs)
	validationError = validationError.Merge(parseCategoryIDsError)

	if validationError != nil && !validationError.IsEmpty() {
		return input{}, validationError
	}

	isPublished, _ := domain.ParseBool(req.IsPublished)
	publishStartAt, _ := domain.ParseOptionalDatetime(req.PublishStartAt)
	publishEndAt, _ := domain.ParseOptionalDatetime(req.PublishEndAt)

	return input{
		Title:          req.Title,
		Body:           req.Body,
		IsPublished:    isPublished,
		PublishStartAt: publishStartAt,
		PublishEndAt:   publishEndAt,
		CategoryIDs:    categoryIDs,
	}, nil
}

func parseCategoryIDs(values []string) ([]category.CategoryID, *handlererror.ValidationError) {
	if len(values) == 0 {
		return []category.CategoryID{}, nil
	}

	categoryIDs := make([]category.CategoryID, 0, len(values))
	for _, rawID := range values {
		if strings.TrimSpace(rawID) == "" {
			continue
		}
		categoryID, err := category.ParseCategoryID(strings.TrimSpace(rawID))
		if err != nil {
			return nil, &handlererror.ValidationError{Messages: handlererror.ValidationErrorMessages{
				"categoryIds": "選択したカテゴリが不正です。",
			}}
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	return categoryIDs, nil
}

func getValidationMessage(field, tag string) string {
	switch field {
	case "title":
		switch tag {
		case "required":
			return "タイトルを入力してください。"
		case "notblank":
			return "タイトルを入力してください。"
		case "max":
			return "タイトルは255文字以下で入力してください。"
		}
	case "body":
		switch tag {
		case "required":
			return "本文を入力してください。"
		case "notblank":
			return "本文を入力してください。"
		}
	case "isPublished":
		switch tag {
		case "required":
			return "公開状態を指定してください。"
		case "bool":
			return "公開状態が不正です。"
		}
	case "publishStartAt":
		switch tag {
		case "datetime":
			return "日時の形式が不正です。"
		case "datetime_lt":
			return "公開開始時刻は公開終了時刻より前を指定してください。"
		}
	case "publishEndAt":
		switch tag {
		case "datetime":
			return "日時の形式が不正です。"
		}
	}
	return "入力内容が不正です。"
}
