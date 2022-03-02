package service

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService struct for rep
type UserService struct {
	conn repository.Users
}

// NewUserServicePostgres used for setting postgres services
func NewUserServicePostgres(rep *repository.PostgresUser) *UserService {
	return &UserService{conn: rep}
}

// NewUserServiceMongo used for setting mongo services
func NewUserServiceMongo(rep *repository.MongoUser) *UserService {
	return &UserService{conn: rep}
}

func NewUserService(r repository.Users) *UserService {
	return &UserService{conn: r}
}

// Add record about user
func (srv *UserService) Add(e context.Context, user *model.User) error {
	hashedPass, ok := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if ok != nil {
		return fmt.Errorf("can't hashing password %w", ok)
	}
	user.Password = string(hashedPass)

	_, err := srv.conn.Get(e, user)
	if err != nil {
		return srv.conn.Add(e, user)
	}
	return fmt.Errorf("this user already exist %w", err)
}

// GetAll getting all users
func (srv *UserService) GetAll(e context.Context) ([]*model.User, error) {
	return srv.conn.GetAll(e)
}

// Get getting parking lot by username
func (srv *UserService) Get(e context.Context, user *model.User) (*model.User, error) {
	getUser, err := srv.conn.Get(e, user)
	if err != nil {
		return nil, err
	}
	if ok := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(user.Password)); ok != nil && user.Password != "" {
		return nil, fmt.Errorf("authentication comparing passwords error %v", ok)
	}
	return getUser, err
}

// Update updating user
func (srv *UserService) Update(e context.Context, user *model.User) error {
	return srv.conn.Update(e, user)
}

// Delete deleting user
func (srv *UserService) Delete(e context.Context, username string) error {
	return srv.conn.Delete(e, username)
}
