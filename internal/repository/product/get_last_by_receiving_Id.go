package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"database/sql"
	"errors"
)

func (r *productRepo) GetLastProductByReceigingId(ctx context.Context, receptionID string) (*entity.Product, error) {
	var product entity.Product

	err := r.db.QueryRow(
		ctx,
		`SELECT id, receiving_id, category_id, added_at 
		FROM products 
		WHERE receiving_id = $1 
		ORDER BY added_at DESC 
		LIMIT 1`,
		receptionID,
	).Scan(&product.ID, &product.ReceptionID, &product.Category, &product.DateTime)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.ErrNoProductsToDelete
		}
		return nil, err
	}
	return &product, nil
}
