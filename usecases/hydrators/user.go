package hydrators

import (
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/protousers"
)

func ProtoUser(user models.User) *protousers.User {
	return &protousers.User{
		Id:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Username:     user.Username,
		Email:        user.Email,
		Bio:          user.Bio,
		ProfileImage: user.ProfileImage,
		CoverImage:   user.CoverImage,
	}
}
