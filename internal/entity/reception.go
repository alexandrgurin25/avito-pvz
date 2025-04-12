package entity

import "time"

type Reception struct {
	ID        string    // Уникальный идентификатор приемки (UUID)
	PvzID     string    // ID ПВЗ, к которому относится приемка
	DateTime  time.Time // Дата и время проведения приемки
	CloseTime time.Time
	Status    string // Статус: "in_progress" или "close"
}
