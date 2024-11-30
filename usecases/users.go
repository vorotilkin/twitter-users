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
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserByID(ctx context.Context, id int32) (models.User, error)
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
		User: toProtoUser(user),
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

func (s *UsersServer) UserByEmail(ctx context.Context, request *proto.UserByEmailRequest) (*proto.UserByEmailResponse, error) {
	user, err := s.usersRepository.UserByEmail(ctx, request.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.ID == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.UserByEmailResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *UsersServer) UserByID(ctx context.Context, request *proto.UserByIDRequest) (*proto.UserByIDResponse, error) {
	user, err := s.usersRepository.UserByID(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.ID == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.UserByIDResponse{
		User: toProtoUser(user),
	}, nil
}

func NewUsersServer(usersRepo UsersRepository) *UsersServer {
	return &UsersServer{
		usersRepository: usersRepo,
	}
}

func toProtoUser(user models.User) *proto.User {
	return &proto.User{
		Id:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Username:     user.Username,
		Email:        user.Email,
	}
}
