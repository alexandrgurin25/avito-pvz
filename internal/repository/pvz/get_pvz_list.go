package pvz

import (
	"avito-pvz/internal/entity"
	"context"
	"fmt"
	"time"
)

func (r *pvzRepository) GetPVZsWithFilters(
	ctx context.Context,
	startDate, endDate *time.Time,
	page, limit int,
) ([]entity.PVZ, error) {

	query := `
	SELECT 
		pvz.id, c.name, pvz.created_at,
		r.id, r.start_time, r.status, r.end_time,
		p.id, p.added_at, cat.name
	FROM pvz
	JOIN cities c ON pvz.city_id = c.id
	JOIN receivings r ON pvz.id = r.pvz_id
	JOIN products p ON r.id = p.receiving_id
	JOIN categories cat ON p.category_id = cat.id
	where pvz.id IN (SELECT pvz.id FROM pvz JOIN receivings ON pvz.id = receivings.pvz_id
	WHERE receivings.start_time >= $1
	  AND receivings.end_time <= $2
	ORDER BY pvz.created_at DESC
	LIMIT $3 OFFSET $4)
	ORDER BY pvz.created_at DESC
    LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, startDate, endDate, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pvz with receptions: %w", err)
	}
	defer rows.Close()

	// Используем мапы для группировки данных
	pvzMap := make(map[string]*entity.PVZ)
	receptionMap := make(map[string]*entity.Reception)

	for rows.Next() {
		var (
			pvzID, cityName, recID, status, prodID, prodType string
			pvzCreated, recStart, recEnd, prodAddedAt        time.Time
		)

		err := rows.Scan(&pvzID, &cityName, &pvzCreated,
			&recID, &recStart, &status, &recEnd,
			&prodID, &prodAddedAt, &prodType,
		)
		if err != nil {
			return nil, err
		}

		// Обработка PVZ
		if _, exists := pvzMap[pvzID]; !exists {
			pvzMap[pvzID] = &entity.PVZ{
				UUID:       pvzID,
				City:       entity.City{Name: cityName},
				CreatedAt:  pvzCreated,
				Receptions: []*entity.Reception{},
			}
		}

		// Обработка Reception
		if _, exists := receptionMap[recID]; !exists {
			reception := &entity.Reception{
				ID:        recID,
				PvzID:     pvzID,
				DateTime:  recStart,
				CloseTime: recEnd,
				Status:    status,
				Products:  []entity.Product{},
			}
			receptionMap[recID] = reception
			pvzMap[pvzID].Receptions = append(pvzMap[pvzID].Receptions, reception)
		}

		// Добавление продукта
		product := entity.Product{
			ID:          prodID,
			ReceptionID: recID,
			DateTime:    prodAddedAt,
			Category:    prodType,
		}
		receptionMap[recID].Products = append(receptionMap[recID].Products, product)
	}

	// Преобразование мапы в срез
	result := make([]entity.PVZ, 0, len(pvzMap))
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	return result, nil
}
