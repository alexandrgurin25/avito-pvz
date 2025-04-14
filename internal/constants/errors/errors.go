package myerrors

import "errors"

var (
	// Ошибки, связанные с ролями и токенами
	ErrInvalidRole = errors.New("недопустимая роль")
	// ErrInvalidGenerateJWT  = errors.New("не удалось сгенерировать токен")
	ErrInvalidOrExpiredJWT = errors.New("недопустимый или истекший токен")
	ErrAccessDenied        = errors.New("доступ запрещен")

	// Ошибки, связанные с JSON
	ErrJsonNotFound = errors.New("JSON не найден")

	// Ошибки, связанные с городами
	ErrCityOrIdNil  = errors.New("ID и город обязательны")
	ErrCityNotFound = errors.New("город не найден")

	// Ошибки, связанные с ПВЗ (пункты выдачи заказов)
	ErrPvzIdNil         = errors.New("pvzId обязательны")
	ErrInvalidCity      = errors.New("ПВЗ может быть создан только в Москве, Санкт-Петербурге или Казани")
	ErrPVZAlreadyExists = errors.New("ПВЗ уже существует")
	ErrPVZNotFound      = errors.New("ПВЗ не найден")

	// Ошибки, связанные с активными приемами
	ErrActiveReceptionNotFound = errors.New("недопустимый запрос или нет активного приема")
	ErrActiveReceptionFound    = errors.New("активный прием найден")

	// Ошибки, связанные с параметрами ПВЗ
	ErrPvzIdOrTypeNil = errors.New("pvzId и тип обязательны")

	// Ошибки, связанные с типами продуктов
	ErrInvalidProductType = errors.New(
		"недопустимый тип продукта: должен быть одним из [электроника, одежда, обувь]",
	)

	// Ошибки, связанные с учетными данными
	ErrEmailOrPasswordEmpty = errors.New("email пользователя и пароль обязательны")
	ErrUserNotFound         = errors.New("пользователь с такими данными не найден")
	ErrWrongPassword        = errors.New("wrong password")
	ErrUserAlreadyExists    = errors.New("пользователь с такими данными уже существует")

	// Ошибки, связанные с продуктами
	ErrNoProductsToDelete = errors.New("нет продуктов для удаления")
)
