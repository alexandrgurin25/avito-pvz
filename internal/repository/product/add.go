package product

import (
	"avito-pvz/internal/entity"
	"context"
)

func (r *productRepo) AddProduct(ctx context.Context,
	receivingId string,
	categoryId int) (*entity.Product, error) {
	var product entity.Product

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO products (receiving_id, category_id) 
		VALUES ($1, $2)
		RETURNING id, receiving_id, category_id, added_at`,
		receivingId,
		categoryId,
	).Scan(&product.ID, &product.ReceptionID, &product.Category, &product.DateTime)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
