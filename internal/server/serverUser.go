// Package server grpc
package server

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/EgMeln/CRUDentity/protocol"
	log "github.com/sirupsen/logrus"
)

// UserServer for grpc
type UserServer struct {
	userService  *service.UserService
	authService  *service.AuthenticationService
	imageService service.ImageStore
	protocol.UnimplementedUserServiceServer
}

// NewUserServer user server
func NewUserServer(userService *service.UserService, authService *service.AuthenticationService, imageService service.ImageStore) *UserServer {
	return &UserServer{userService: userService, authService: authService, imageService: imageService}
}

// AddUser sign up user
func (srv *UserServer) AddUser(ctx context.Context, in *protocol.SignUpUserRequest) (*protocol.SignUpUserResponse, error) {
	user := model.User{Username: in.Username, Password: in.Password}
	err := srv.userService.Add(ctx, &user)
	if err != nil {
		log.Warnf("grpc add user: %v", err)
		return nil, fmt.Errorf("grpc add user %w", err)
	}
	res := &protocol.SignUpUserResponse{Username: in.Username, Password: in.Password}
	return res, nil
}

// GetUser getting user by id
func (srv *UserServer) GetUser(ctx context.Context, in *protocol.GetUserRequest) (*protocol.GetUserResponse, error) {
	user, err := srv.userService.Get(ctx, &model.User{Username: in.Username})
	if err != nil {
		log.Warnf("grpc get user: %v", err)
		return nil, fmt.Errorf("grpc get user %w", err)
	}
	res := &protocol.GetUserResponse{Username: user.Username, Password: user.Password}
	return res, nil
}

// GetAllUser getting all users
func (srv *UserServer) GetAllUser(ctx context.Context, in *protocol.Empty) (*protocol.GetAllUsersResponse, error) {
	users, err := srv.userService.GetAll(ctx)
	if err != nil {
		log.Warnf("grpc get all users: %v", err)
		return nil, fmt.Errorf("grpc get all user %w", err)
	}
	var resUsers []*protocol.User

	for i := 0; i < len(users); i++ {
		resUsers[i] = &protocol.User{Username: users[i].Username, Password: users[i].Password}
	}
	log.Info(resUsers)
	res := &protocol.GetAllUsersResponse{User: resUsers}
	log.Info(res)
	return res, nil
}

// UpdateUser updating user
func (srv *UserServer) UpdateUser(ctx context.Context, in *protocol.UpdateUserRequest) (*protocol.UpdateUserResponse, error) {
	err := srv.userService.Update(ctx, &model.User{Username: in.Username, Password: in.Password})
	if err != nil {
		log.Warnf("grpc update user: %v", err)
		return nil, fmt.Errorf("grpc update user %w", err)
	}
	res := &protocol.UpdateUserResponse{Username: in.Username, Password: in.Password}
	return res, nil
}

// DeleteUser deleting user
func (srv *UserServer) DeleteUser(ctx context.Context, in *protocol.DeleteUserRequest) (*protocol.Empty, error) {
	err := srv.userService.Delete(ctx, in.Username)
	if err != nil {
		log.Warnf("grpc delete user: %v", err)
		return nil, fmt.Errorf("grpc delete user %w", err)
	}
	return &protocol.Empty{}, nil
}

// SignIn authentication user
func (srv *UserServer) SignIn(ctx context.Context, in *protocol.SignInUserRequest) (*protocol.Tokens, error) {
	user, err := srv.userService.Get(ctx, &model.User{Username: in.Username, Password: in.Password})
	if err != nil {
		log.Warnf("grpc authenticate: %v", err)
		return nil, fmt.Errorf("grpc authenticate %w", err)
	}
	accessToken, refreshToken, err := srv.authService.SignIn(ctx, user)
	log.Infof(accessToken, refreshToken, err)
	if err != nil {
		log.Warnf("grpc sign in : %v", err)
		return nil, fmt.Errorf("grpc sign in %w", err)
	}
	res := &protocol.Tokens{Access: accessToken, Refresh: refreshToken}
	return res, nil
}

// RefreshToken refreshing tokens
func (srv *UserServer) RefreshToken(ctx context.Context, in *protocol.RefreshTokenRequest) (*protocol.Tokens, error) {
	user := &model.User{Username: in.Username}
	accessToken, refreshToken, err := srv.authService.SignIn(ctx, user)
	if err != nil {
		log.Warnf("grpc sign in : %v", err)
		return nil, fmt.Errorf("grpc sign in %w", err)
	}
	res := &protocol.Tokens{Access: accessToken, Refresh: refreshToken}
	return res, nil
}

// UploadImage image
func (srv *UserServer) UploadImage(stream protocol.UserService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("cannot receive image info %v", err)
	}
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request with image type %s", imageType)
	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		log.Print("waiting to receive more data")
		req, ok := stream.Recv()
		if ok == io.EOF {
			log.Print("no more data")
			break
		}
		if ok != nil {
			return fmt.Errorf("cannot receive chunk data: %v", err)
		}
		chunk := req.GetChunkData()
		size := len(chunk)
		log.Printf("received a chunk with size: %d", size)
		imageSize += size
		if imageSize > 1<<20 {
			return fmt.Errorf("image is too large: %d > %d", imageSize, 1<<20)
		}
		_, ok = imageData.Write(chunk)
		if ok != nil {
			return fmt.Errorf("cannot write chunk data: %v", err)
		}
	}
	imageID, err := srv.imageService.Save(imageType, imageData)
	if err != nil {
		return fmt.Errorf("cannot save image : %v", err)
	}
	res := &protocol.UploadImageResponse{
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return fmt.Errorf("cannot send response: %v", err)
	}

	log.Printf("saved image with id: %s, size: %d", imageID, imageSize)
	return nil
}
