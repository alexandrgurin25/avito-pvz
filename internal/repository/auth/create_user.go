package auth

import (
	"avito-pvz/internal/entity"
	"context"
)

func (r *authRepo) CreateUser(
	ctx context.Context, email string,
	passwordHash string,
	role string) (*entity.User, error) {

	var user entity.User

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO users(email, password_hash, role)
		VALUES($1, $2, $3)
		RETURNING id, email, password_hash, role`,
		email,
		passwordHash,
		role,
	).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
