package usecases

import (
	"context"
	"twitter-users/domain/models"
	"twitter-users/proto"
)

type UsersRepository interface {
	Create(ctx context.Context, name, passwordHash, username, email string) (models.User, error)
}

type UsersServer struct {
	proto.UnimplementedUsersServer
	usersRepository UsersRepository
}

func (s *UsersServer) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	user, err := s.usersRepository.Create(ctx, request.GetName(), request.GetPasswordHash(), request.GetUsername(), request.GetEmail())
	if err != nil {
		return nil, err
	}

	return &proto.CreateResponse{
		Id:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

func NewUsersServer(usersRepo UsersRepository) *UsersServer {
	return &UsersServer{
		usersRepository: usersRepo,
	}
}
