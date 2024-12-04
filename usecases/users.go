package usecases

import (
	"context"
	"errors"
	"github.com/samber/lo"
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
	UsersByIDs(ctx context.Context, ids []int32) ([]models.User, error)
	UpdateByID(ctx context.Context, userToUpdate models.UserOption) (bool, error)
	Follow(ctx context.Context, userID, targetUserID int32) (bool, error)
	Unfollow(ctx context.Context, userID, targetUserID int32) (bool, error)
	NewUsers(ctx context.Context, limit int32) ([]models.User, error)
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

func (s *UsersServer) UsersByIDs(ctx context.Context, request *proto.UsersByIDsRequest) (*proto.UsersByIDsResponse, error) {
	userIDs := request.GetIds()
	if len(userIDs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no ids provided")
	}

	users, err := s.usersRepository.UsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UsersByIDsResponse{
		Users: hydrators.ProtoUsers(users),
	}, nil
}

func (s *UsersServer) UpdateByID(ctx context.Context, request *proto.UpdateByIDRequest) (*proto.UpdateByIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	userID := request.GetId()
	if userID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	userToUpdate := models.UserOption{
		ID:           userID,
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

	users, err := s.usersRepository.UsersByIDs(ctx, []int32{userID})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user, exist := lo.Find(users, func(user models.User) bool {
		return user.ID == userID
	})

	if !exist {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.UpdateByIDResponse{
		User: hydrators.ProtoUser(user),
	}, nil
}

func (s *UsersServer) Follow(ctx context.Context, request *proto.FollowRequest) (*proto.FollowResponse, error) {
	userID := request.GetUserId()
	if userID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	targetUserID := request.GetTargetUserId()
	if targetUserID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid target id")
	}

	if userID == targetUserID {
		return nil, status.Error(codes.InvalidArgument, "ids equals")
	}

	var (
		ok  bool
		err error
	)

	switch request.GetOperationType() {
	case proto.FollowRequest_OPERATION_TYPE_UNFOLLOW:
		ok, err = s.usersRepository.Unfollow(ctx, userID, targetUserID)
	default:
		ok, err = s.usersRepository.Follow(ctx, userID, targetUserID)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.FollowResponse{Ok: ok}, nil
}

func (s *UsersServer) NewUsers(ctx context.Context, request *proto.NewUsersRequest) (*proto.NewUsersResponse, error) {
	users, err := s.usersRepository.NewUsers(ctx, request.GetLimit())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.NewUsersResponse{Users: hydrators.ProtoUsers(users)}, nil
}

func NewUsersServer(usersRepo UsersRepository) *UsersServer {
	return &UsersServer{
		usersRepository: usersRepo,
	}
}
