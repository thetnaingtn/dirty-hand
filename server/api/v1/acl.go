package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thetnaingtn/dirty-hand/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextKey int

const (
	userIdContextKey ContextKey = iota
	sessionIdContextKey
)

var authticationAllowListMethods = map[string]bool{
	"/api.v1.UserService/CreateUser":    true,
	"/api.v1.UserService/CreateSession": true,
}

type GRPCAuthInterceptor struct {
	store *store.Store
}

func NewGRPCAuthInterceptor(store *store.Store) *GRPCAuthInterceptor {
	return &GRPCAuthInterceptor{
		store: store,
	}
}

func (in *GRPCAuthInterceptor) AuthenticateInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "failed to parse metadata")
	}

	if sessionCookieValue, err := getSessionIDFromMetadata(md); err != nil && sessionCookieValue != "" {
		user, err := in.authenticateBySession(ctx, sessionCookieValue)
		if err == nil && user != nil {
			_, sessionId, parsedErr := parseSessionCookieValue(sessionCookieValue)
			if parsedErr != nil {
				return nil, status.Error(codes.Unauthenticated, "failed to parsed session cookie")
			}

			ctx = context.WithValue(ctx, userIdContextKey, user.ID)
			if sessionId != "" {
				ctx = context.WithValue(ctx, sessionIdContextKey, sessionId)
			}

			if err := in.updateLastAccessedTime(ctx, sessionId); err != nil {
				return nil, status.Error(codes.Internal, "failed to update last accessed time")
			}

			handler(ctx, req)
		}

	}

	if isUnauthorizeAllowMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	return nil, status.Error(codes.Unauthenticated, "authentication required")
}

func (in *GRPCAuthInterceptor) updateLastAccessedTime(ctx context.Context, sessionId string) error {
	return in.store.UpdateLastAccessedTime(ctx, sessionId, time.Now())
}

func (in *GRPCAuthInterceptor) authenticateBySession(ctx context.Context, sessionCookieValue string) (*store.User, error) {
	userId, sessionId, err := parseSessionCookieValue(sessionCookieValue)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid session cookie format")
	}

	user, err := in.store.GetUser(ctx, &store.FindUser{
		ID: &userId,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user")
	}

	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	sessions, err := in.store.GetUserSessions(ctx, userId)

	if valid := in.validateUserSession(sessions, sessionId); !valid {
		return nil, status.Error(codes.Unauthenticated, "session invalid or expired")
	}

	return user, nil
}

func (in *GRPCAuthInterceptor) validateUserSession(sessions []store.Session, sessionId string) bool {
	for _, session := range sessions {
		if session.SessionID == sessionId {
			expiredTime := session.LastAccessedTime.Add(14 * 24 * time.Hour)

			if expiredTime.Before(time.Now()) {
				return false
			}
		}
	}

	return true
}

func parseSessionCookieValue(sessionId string) (int64, string, error) {
	splits := strings.SplitN(sessionId, "-", 2)
	if len(splits) > 2 {
		return 0, "", errors.New("invalid session cookie format")
	}

	userId, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid user id in cookie: %v", err)
	}

	return userId, splits[1], nil
}

func getSessionIDFromMetadata(md metadata.MD) (string, error) {
	var sessionId string

	for _, t := range append(md.Get("grpcgateway-cookie"), md.Get("Cookie")...) {
		header := http.Header{}
		header.Add("Cookie", t)

		request := http.Request{Header: header}

		if v, _ := request.Cookie("user_session"); v != nil {
			sessionId = v.Value
		}
	}

	if sessionId == "" {
		return "", errors.New("session cookie not found")
	}

	return sessionId, nil
}

func isUnauthorizeAllowMethod(method string) bool {
	return authticationAllowListMethods[method]
}
