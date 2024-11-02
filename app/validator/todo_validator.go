package validator

import (
	"app/graph/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateTodo(input model.CreateTodoInput) error {
	return validation.ValidateStruct(&input,
		validation.Field(
			&input.Title,
			validation.Required.Error("タイトルは必須入力です。"),
			validation.RuneLength(1, 50).Error("タイトルは1 ~ 50文字での入力をお願いします。"),
		),
	)
}

func ValidateUpdateTodo(input model.UpdateTodoInput) error {
	return validation.ValidateStruct(&input,
		validation.Field(
			&input.Title,
			validation.Required.Error("タイトルは必須入力です。"),
			validation.RuneLength(1, 50).Error("タイトルは1 ~ 50文字での入力をお願いします。"),
		),
	)
}
