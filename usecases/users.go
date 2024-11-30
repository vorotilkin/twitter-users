package usecases

import (
	"context"
	"errors"
	"github.com/samber/mo"
	"github.com/vorotilkin/twitter-users/domain/models"
	"github.com/vorotilkin/twitter-users/proto"
	"github.com/vorotilkin/twitter-users/usecases/hydrators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersRepository interface {
	Create(ctx context.Context, name, passwordHash, username, email string) (models.User, error)
	FetchPasswordHashByEmail(ctx context.Context, email string) (string, error)
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserByID(ctx context.Context, id int32) (models.User, error)
	UpdateByID(ctx context.Context, userToUpdate models.UserOption) (bool, error)
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
		User: hydrators.ProtoUser(user),
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
		User: hydrators.ProtoUser(user),
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
		User: hydrators.ProtoUser(user),
	}, nil
}

func (s *UsersServer) UpdateByID(ctx context.Context, request *proto.UpdateByIDRequest) (*proto.UpdateByIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if request.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	userToUpdate := models.UserOption{
		ID:           request.GetId(),
		Name:         mo.PointerToOption(request.Name),
		Username:     mo.PointerToOption(request.Username),
		Bio:          mo.PointerToOption(request.Bio),
		ProfileImage: mo.PointerToOption(request.ProfileImage),
		CoverImage:   mo.PointerToOption(request.CoverImage),
	}

	ok, err := s.usersRepository.UpdateByID(ctx, userToUpdate)
	if errors.Is(err, models.ErrNothingToUpdate) {
		return nil, status.Error(codes.InvalidArgument, "nothing to update")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	user, err := s.usersRepository.UserByID(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.ID == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.UpdateByIDResponse{
		User: hydrators.ProtoUser(user),
	}, nil
}

func NewUsersServer(usersRepo UsersRepository) *UsersServer {
	return &UsersServer{
		usersRepository: usersRepo,
	}
}
