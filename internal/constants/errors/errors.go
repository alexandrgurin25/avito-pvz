package myerrors

import "errors"

var (
	ErrInvalidRole         = errors.New("invalid role")
	ErrInvalidGenerateJWT  = errors.New("failed to generate token")
	ErrInvalidOrExpiredJWT = errors.New("invalid or expired token")
	ErrAccessDenied        = errors.New("access denied")
)

var (
	ErrCityOrIdNil      = errors.New("ID and City are required")
	ErrInvalidCity      = errors.New("PVZ can only be created in Moscow, Saint Petersburg or Kazan")
	ErrPVZAlreadyExists = errors.New("PVZ already exists")

	ErrCityNotFound = errors.New("city not found")
)
