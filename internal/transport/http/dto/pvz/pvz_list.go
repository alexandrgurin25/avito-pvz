package pvz

import (
	"avito-pvz/internal/transport/http/dto/product"
	"avito-pvz/internal/transport/http/dto/reception"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type GetPVZsRequest struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02T15:04:05Z"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02T15:04:05Z"`
	Page      int        `form:"page,default=1"`
	Limit     int        `form:"limit,default=10"`
}

type PVZWithReceptions struct {
	PVZ        PvzResponse             `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions,omitempty"`
}

type ReceptionWithProducts struct {
	Reception reception.ReceptionResponse `json:"reception"`
	Products  []product.ProductResponse   `json:"products,omitempty"`
}

func (r *GetPVZsRequest) ParseFromQuery(values url.Values) error {

	if startDateStr := values.Get("startDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			r.StartDate = &t
		} else {
			return fmt.Errorf("invalid startDate format")
		}
	}

	if endDateStr := values.Get("endDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			r.EndDate = &t
		} else {
			return fmt.Errorf("invalid endDate format")
		}
	}

	if pageStr := values.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			r.Page = page
		}
	}

	if limitStr := values.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			r.Limit = limit
		}
	}

	return nil
}
