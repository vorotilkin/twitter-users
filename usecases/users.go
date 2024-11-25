package usecases

import (
	"context"
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersRepository interface {
	Create(ctx context.Context, name, passwordHash, username, email string) (models.User, error)
	FetchPasswordHashByEmail(ctx context.Context, email string) (string, error)
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

func (s *UsersServer) PasswordHashByEmail(ctx context.Context, request *proto.PasswordHashByEmailRequest) (*proto.PasswordHashByEmailResponse, error) {
	hash, err := s.usersRepository.FetchPasswordHashByEmail(ctx, request.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(hash) == 0 {
		return nil, status.Error(codes.NotFound, "hash not found")
	}

	return &proto.PasswordHashByEmailResponse{PasswordHash: hash}, nil
}

func NewUsersServer(usersRepo UsersRepository) *UsersServer {
	return &UsersServer{
		usersRepository: usersRepo,
	}
}
