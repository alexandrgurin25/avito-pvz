package entity

import "time"

type Product struct {
	ID          string    // Уникальный идентификатор товара (UUID)
	ReceptionID string    // ID приемки, к которой относится товар
	DateTime    time.Time // Дата и время добавления товара
	Category    string    // Тип товара: "электроника", "одежда" или "обувь"
}
