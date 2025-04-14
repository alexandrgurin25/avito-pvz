package auth

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"database/sql"
	"errors"
)

func (r *authRepo) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	err := r.db.QueryRow(
		ctx,
		`SELECT id, email, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&user.Id, &user.Email, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
