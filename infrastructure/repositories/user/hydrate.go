package user

import (
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/schema/gen/my_database/public/model"
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
