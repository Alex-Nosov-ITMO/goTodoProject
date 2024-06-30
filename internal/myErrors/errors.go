package myErrors

import (
	"fmt"
)

var (
	BadRequest = "неверный запрос"
	InternalServerError = "ошибка сервера"
	Unauthorized = "не авторизован"
	NotFound = "не найдено"
	InvalidID = "невалидный идентификатор"
	WrongPassword = "неверный пароль, не авторизован"
)

func WithMassage(message string, err error) error {
	return fmt.Errorf("%s: %w", message, err)
}