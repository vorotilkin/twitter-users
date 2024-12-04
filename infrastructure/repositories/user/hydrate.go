package user

import (
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/schema/gen/my_database/public/model"
)

func toDomain(user model.User, followingIDs []int32, followerIDs []int32) models.User {
	return models.User{
		ID:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Username:     user.Username,
		Email:        user.Email,
		Bio:          lo.FromPtr(user.Bio),
		ProfileImage: lo.FromPtr(user.ProfileImage),
		CoverImage:   lo.FromPtr(user.CoverImage),
		FollowingIDs: followingIDs,
		FollowerIDs:  followerIDs,
	}
}
