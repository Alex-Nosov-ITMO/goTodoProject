package myErrors

import (
	"errors"
)

var ErrorsString = map[string] string{
	"BadRequest": errors.New("неверный запрос").Error(),
	"InternalServerError": errors.New("ошибка сервера").Error(),
	"Unauthorized": errors.New("не авторизован").Error(),
	"NotFound": errors.New("не найдено").Error(),
	"InvalidID": errors.New("невалидный идентификатор").Error(),
	"WrongPassword": errors.New("неверный пароль, не авторизован").Error(),
}

func WithMassage(message string, err error) error {
	return errors.New(message + ": " + err.Error())
}