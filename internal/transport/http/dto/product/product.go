package product

import "time"

type AddProductRequest struct {
	Type  string `json:"type"`
	PvzId string `json:"pvzId"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}
