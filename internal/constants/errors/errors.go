package myerrors

import "errors"

var (
	ErrInvalidRole         = errors.New("invalid role")
	ErrInvalidGenerateJWT  = errors.New("failed to generate token")
	ErrInvalidOrExpiredJWT = errors.New("invalid or expired token")
	ErrAccessDenied        = errors.New("access denied")
)

var (
	ErrJsonNotFound = errors.New("JSON was not found")

	ErrCityOrIdNil  = errors.New("ID and City are required")
	ErrCityNotFound = errors.New("city not found")

	ErrPvzIdNil         = errors.New("pvzId are required")
	ErrInvalidCity      = errors.New("PVZ can only be created in Москва, Санкт-Петербург or Казань")
	ErrPVZAlreadyExists = errors.New("PVZ already exists")
	ErrPVZNotFound      = errors.New("PVZ not found")

	ErrActiveReceptionNotFound = errors.New("invalid request or there is no active reception")
	ErrActiveReceptionFound    = errors.New("active reception was found")

	ErrPvzIdOrTypeNil = errors.New("pvzId and type are required")

	ErrInvalidProductType = errors.New(
		"invalid product type: must be one of [электроника, одежда, обувь]",
	)
)
