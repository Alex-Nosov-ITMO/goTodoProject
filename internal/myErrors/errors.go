package myErrors

import (
	"fmt"
)
/*
var ErrorsString = map[string] string{
	"BadRequest": errors.New("неверный запрос").Error(),
	"InternalServerError": errors.New("ошибка сервера").Error(),
	"Unauthorized": errors.New("не авторизован").Error(),
	"NotFound": errors.New("не найдено").Error(),
	"InvalidID": errors.New("невалидный идентификатор").Error(),
	"WrongPassword": errors.New("неверный пароль, не авторизован").Error(),
}
*/
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