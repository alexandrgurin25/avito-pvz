package entity

type User struct {
	Id           int
	Email        string
	PasswordHash string
	Role         string
}
