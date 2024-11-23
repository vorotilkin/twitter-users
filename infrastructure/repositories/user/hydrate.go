package user

import (
	"twitter-users/domain/models"
	"twitter-users/schema/gen/my_database/public/model"
)

func toDomain(user model.User) models.User {
	return models.User{
		ID:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Username:     user.Username,
		Email:        user.Email,
	}
}
