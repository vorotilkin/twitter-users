package hydrators

import (
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/proto"
)

func ProtoUsers(users []models.User) []*proto.User {
	return lo.Map(users, func(user models.User, _ int) *proto.User {
		return ProtoUser(user)
	})
}

func ProtoUser(user models.User) *proto.User {
	return &proto.User{
		Id:               user.ID,
		Name:             user.Name,
		PasswordHash:     user.PasswordHash,
		Username:         user.Username,
		Email:            user.Email,
		Bio:              user.Bio,
		ProfileImage:     user.ProfileImage,
		CoverImage:       user.CoverImage,
		FollowingUserIds: user.FollowingIDs,
	}
}
