package service

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MyMockedUser struct {
	user mocks.Users
}

func TestUserService_Get(t *testing.T) {
	testUser := new(MyMockedUser)
	user := &model.User{Username: "test", Password: "", Admin: false}
	testUser.user.On("Get", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(
		func(e context.Context, user *model.User) *model.User {
			return &model.User{Username: "test", Password: "", Admin: false}
		},
		func(e context.Context, user *model.User) error {
			return nil
		})
	serviceUser := NewUserService(&testUser.user)
	_, err := serviceUser.Get(context.Background(), user)
	require.NoError(t, err)
}
func TestUserService_Add(t *testing.T) {
	testUser := new(MyMockedUser)
	user := &model.User{Username: "test", Password: "1234", Admin: false}
	testUser.user.On("Get", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(
		func(e context.Context, user *model.User) *model.User {
			return &model.User{Username: "test", Password: "1234", Admin: false}
		},
		func(e context.Context, user *model.User) error {
			return fmt.Errorf("this user already exist")
		})
	testUser.user.On("Add", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(
		func(e context.Context, user *model.User) error {
			return nil
		})
	serviceUser := NewUserService(&testUser.user)
	err := serviceUser.Add(context.Background(), user)
	require.NoError(t, err)
}
func TestUserService_GetAll(t *testing.T) {
	testUser := new(MyMockedUser)
	testUser.user.On("GetAll", mock.AnythingOfType("*context.emptyCtx")).Return(
		func(e context.Context) []*model.User {
			return []*model.User{}
		},
		func(e context.Context) error {
			return nil
		})
	serviceUser := NewUserService(&testUser.user)
	_, err := serviceUser.GetAll(context.Background())
	require.NoError(t, err)
}
func TestUserService_Delete(t *testing.T) {
	testUser := new(MyMockedUser)
	testUser.user.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(
		func(e context.Context, username string) error {
			return nil
		})
	serviceUser := NewUserService(&testUser.user)
	username := "test"
	err := serviceUser.Delete(context.Background(), username)
	require.NoError(t, err)
}
func TestUserService_Update(t *testing.T) {
	testUser := new(MyMockedUser)
	testUser.user.On("Update", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(
		func(e context.Context, user *model.User) error {
			return nil
		})
	serviceUser := NewUserService(&testUser.user)
	user := &model.User{Username: "test", Password: "1234", Admin: true}
	err := serviceUser.Update(context.Background(), user)
	require.NoError(t, err)
}
