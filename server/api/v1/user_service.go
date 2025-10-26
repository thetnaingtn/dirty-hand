package v1

import (
	"context"

	apiv1 "github.com/thetnaingtn/dirty-hand/proto/gen/api/v1"
	"github.com/thetnaingtn/dirty-hand/store"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *APIV1Service) CreateUser(ctx context.Context, req *apiv1.CreateUserRequest) (*apiv1.User, error) {
	adminRoleType := store.RoleAdmin
	existingUsers, err := s.store.ListUsers(ctx, &store.FindUser{Role: &adminRoleType})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing admin users: %v", err)
	}

	user := &store.User{
		Username: req.GetUsername(),
		Role:     store.RoleProductView,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user.PasswordHash = string(hashedPassword)

	roleToAssign := store.RoleAdmin
	if len(existingUsers) > 0 {
		roleToAssign = store.RoleProductView
	} else {
		roleToAssign = store.RoleAdmin
	}

	user.Role = roleToAssign

	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &apiv1.User{
		Id:       user.ID,
		Username: user.Username,
		Role:     s.convertUserRoleFromStore(user.Role),
		Password: user.PasswordHash,
	}, nil
}

func (s *APIV1Service) CreateSession(ctx context.Context, req *apiv1.CreateSessionRequest) (*apiv1.CreateSessionResponse, error) {
	return nil, nil
}

func (s *APIV1Service) DeleteSession(ctx context.Context, req *apiv1.DeleteSessionRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) convertUserRoleFromStore(role store.Role) apiv1.Role {
	switch role {
	case store.RoleAdmin:
		return apiv1.Role_ADMIN
	case store.RoleProductEdit:
		return apiv1.Role_PRODUCT_EDITOR
	case store.RoleProductView:
		return apiv1.Role_PRODUCT_VIEWER
	default:
		return apiv1.Role_UNSPECIFIED
	}
}
