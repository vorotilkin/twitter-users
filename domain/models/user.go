package models

type User struct {
	ID           int32
	Name         string
	PasswordHash string
	Username     string
	Email        string
}
