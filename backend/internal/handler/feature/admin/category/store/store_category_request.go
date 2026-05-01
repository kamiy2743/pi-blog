package store

import (
	"net/http"

	"blog/internal/handler/handlererror"
	"blog/internal/handler/handlerrequest"
	"blog/internal/handler/validator"
)

type request struct {
	Name string `json:"name" validate:"required,notblank,max=64"`
}

type input struct {
	Name string
}

func toInput(r *http.Request) (input, *handlererror.ValidationError) {
	req := request{}
	handlerrequest.DecodeJSONForm(r, &req)
	validationError := validator.Validate(req, getValidationMessage)

	return input{
		Name: req.Name,
	}, validationError
}

func getValidationMessage(field, tag string) string {
	switch field {
	case "name":
		switch tag {
		case "required":
			return "カテゴリ名を入力してください。"
		case "notblank":
			return "カテゴリ名を入力してください。"
		case "max":
			return "カテゴリ名は64文字以下で入力してください。"
		}
	}
	return "入力内容が不正です。"
}
