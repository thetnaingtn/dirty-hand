package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	apiv1 "github.com/thetnaingtn/dirty-hand/proto/gen/api/v1"
	"github.com/thetnaingtn/dirty-hand/store"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *APIV1Service) CreateUser(ctx context.Context, req *apiv1.CreateUserRequest) (*apiv1.User, error) {
	adminRoleType := store.RoleAdmin
	existingUsers, err := s.store.ListUsers(ctx, &store.FindUser{Role: &adminRoleType})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing admin users: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	roleToAssign := store.RoleAdmin
	if len(existingUsers) > 0 {
		roleToAssign = store.RoleProductView
	} else {
		roleToAssign = store.RoleAdmin
	}

	user, err := s.store.CreateUser(ctx, &store.User{
		Username:     req.GetUsername(),
		PasswordHash: string(hashedPassword),
		Role:         roleToAssign,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &apiv1.User{
		Id:       user.ID,
		Username: user.Username,
		Role:     convertUserRoleFromStore(user.Role),
	}, nil
}

func (s *APIV1Service) CreateSession(ctx context.Context, req *apiv1.CreateSessionRequest) (*apiv1.CreateSessionResponse, error) {
	user, err := s.store.GetUser(ctx, &store.FindUser{
		Username: &req.Username,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unmatched username and password")
	}

	lastAccessedAt := time.Now()
	expireTime := time.Now().Add(100 * 365 * 24 * time.Hour)

	if err := s.doSignIn(ctx, user.ID, lastAccessedAt, expireTime); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign in, error: %v", err)
	}

	return &apiv1.CreateSessionResponse{
		User:           convertUserFromStore(user),
		LastAccessedAt: timestamppb.New(lastAccessedAt),
	}, nil
}

func (s *APIV1Service) doSignIn(ctx context.Context, userId int64, lastAccessedAt, expireTime time.Time) error {
	sessionId := uuid.New()
	sessionCookieValue := fmt.Sprintf("%d-%s", userId, sessionId)

	attrs := []string{
		fmt.Sprintf("%s=%s", "user_session", sessionCookieValue),
		"Path=/",
		"HttpOnly",
	}

	if expireTime.IsZero() {
		attrs = append(attrs, "Expires=Thu, 01 Jan 1970 00:00:00 GMT")
	} else {
		attrs = append(attrs, expireTime.Format(time.RFC1123))
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("failed to get metadata from context")
	}
	var origin string
	for _, v := range md.Get("origin") {
		origin = v
	}
	isHTTPS := strings.HasPrefix(origin, "https://")
	if isHTTPS {
		attrs = append(attrs, "SameSite=None")
		attrs = append(attrs, "Secure")
	} else {
		attrs = append(attrs, "SameSite=Strict")
	}

	err := s.store.CreateSession(ctx, &store.Session{
		UserID:           userId,
		SessionID:        sessionId.String(),
		LastAccessedTime: lastAccessedAt,
	})

	if err != nil {
		return status.Error(codes.Internal, "failed to create user session")
	}

	if err := grpc.SetHeader(ctx, metadata.New(map[string]string{
		"Set-Cookie": strings.Join(attrs, ";"),
	})); err != nil {
		return status.Error(codes.Internal, "failed to set grpc cookie")
	}

	return nil
}

func (s *APIV1Service) DeleteSession(ctx context.Context, req *apiv1.DeleteSessionRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func convertUserRoleFromStore(role store.Role) apiv1.Role {
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

func convertUserFromStore(user *store.User) *apiv1.User {
	return &apiv1.User{
		Id:       user.ID,
		Username: user.Username,
		Role:     convertUserRoleFromStore(user.Role),
	}
}
