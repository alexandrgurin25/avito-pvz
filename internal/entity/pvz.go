package entity

import "time"

type City struct {
	Id   int
	Name string
}

type PVZ struct {
	UUID      string
	City      City
	CreatedAt time.Time
}
