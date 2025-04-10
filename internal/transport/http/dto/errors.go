package errors

import "errors"

var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidGenerateJWT = errors.New("failed to generate token")
)
