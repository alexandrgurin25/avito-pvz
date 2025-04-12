package product

import (
	"context"
)

func (r *productRepo) DeleteLastProduct(ctx context.Context, productId string) error {

	_, err := r.db.Exec(
		ctx,
		`DELETE FROM products WHERE id = $1`,
		productId,
	)

	if err != nil {
		return err
	}
	return nil
}
