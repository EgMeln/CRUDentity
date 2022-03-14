package server

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	access  *service.JWTService
	refresh *service.JWTService
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(access, refresh *service.JWTService) *AuthInterceptor {
	return &AuthInterceptor{access: access, refresh: refresh}
}
func (inter *AuthInterceptor) Authorize(ctx context.Context, method string) error {
	if method == "/protobuf.UserService/SignIn" || method == "/protobuf.UserService/AddUser" || method == "/protobuf.UserService/UploadImage" {
		return nil
	}
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}
	tokens := meta["authorization"][0]
	if tokens == "" {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not provided")
	}
	var accessToken string
	_, err := fmt.Sscanf(tokens, "Bearer %s", &accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	return nil
}

// UnaryServerAuthInterceptor returns a server interceptor function to authenticate and authorize unary RPC
func (inter *AuthInterceptor) UnaryServerAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("unary interceptor: ", info.FullMethod)
		err := inter.Authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamServerAuthInterceptor returns a server interceptor function to authenticate and authorize stream RPC
func (inter *AuthInterceptor) StreamServerAuthInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("stream interceptor: ", info.FullMethod)
		err := inter.Authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}
