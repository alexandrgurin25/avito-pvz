package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
)


func (r *pvzRepository) GetCityIdByName(ctx context.Context, city *entity.City) (int, error) {
	var cityID int
	err := r.pool.QueryRow(
		ctx,
		`SELECT id FROM cities WHERE name = $1`,
		city.Name,
	).Scan(&cityID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.ErrCityNotFound
		}
		return 0, fmt.Errorf("failed to get city ID: %v", err)
	}
	return cityID, nil
}
