package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	"context"
	"database/sql"
	"errors"
)

func (r *productRepo) GetIdCategoryByName(ctx context.Context, name string) (int, error) {
	var idCategory int

	err := r.db.QueryRow(
		ctx,
		`SELECT id FROM categories WHERE name = $1`,
		name,
	).Scan(&idCategory)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.ErrInvalidProductType
		}
		return 0, err
	}

	return idCategory, nil
}
