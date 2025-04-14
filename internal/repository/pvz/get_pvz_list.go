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
	`

	rows, err := r.pool.Query(ctx, query, startDate, endDate, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pvz with receptions: %w", err)
	}
	defer rows.Close()

	result := make(map[string]*entity.PVZ)

	for rows.Next() {
		var (
			pvzID, cityName  string
			pvzCreated       time.Time
			recID, status    *string
			recStart, recEnd *time.Time
			prodID, prodType *string
			prodAddedAt      *time.Time
		)

		err := rows.Scan(&pvzID, &cityName, &pvzCreated,
			&recID, &recStart, &status, &recEnd,
			&prodID, &prodAddedAt, &prodType,
		)
		if err != nil {
			return nil, err
		}

		pvz, ok := result[pvzID]
		if !ok {
			pvz = &entity.PVZ{
				UUID:      pvzID,
				City:      entity.City{Name: cityName},
				CreatedAt: pvzCreated,
			}
			result[pvzID] = pvz
		}

		if recID != nil {
			found := false
			for _, r := range pvz.Receptions {
				if r.ID == *recID {
					found = true
					if prodID != nil {
						r.Products = append(r.Products, entity.Product{
							ID:          *prodID,
							DateTime:    *prodAddedAt,
							Category:    *prodType,
							ReceptionID: *recID,
						})
					}
					break
				}
			}
			if !found {
				reception := &entity.Reception{
					ID:       *recID,
					PvzID:    pvzID,
					DateTime: *recStart,
					Status:   *status,
				}
				if recEnd != nil {
					reception.CloseTime = *recEnd
				}
				if prodID != nil {
					reception.CloseTime = *recEnd
					reception.Products = append(reception.Products, entity.Product{
						ID:          *prodID,
						DateTime:    *prodAddedAt,
						Category:    *prodType,
						ReceptionID: *recID,
					})
				}
				pvz.Receptions = append(pvz.Receptions, reception)
			}
		}
	}

	var pvzList []entity.PVZ
	for _, v := range result {
		pvzList = append(pvzList, *v)
	}

	return pvzList, nil
}
