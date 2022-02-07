package service

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
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

// Add record about user
func (srv *UserService) Add(e context.Context, user *model.User) error {
	return srv.conn.Add(e, user)
}

// GetAll getting all users
func (srv *UserService) GetAll(e context.Context) ([]*model.User, error) {
	return srv.conn.GetAll(e)
}

// Get getting parking lot by username
func (srv *UserService) Get(e context.Context, username string) (*model.User, error) {
	return srv.conn.Get(e, username)
}

// Update updating user
func (srv *UserService) Update(e context.Context, username, password string, admin bool) error {
	return srv.conn.Update(e, username, password, admin)
}

// Delete deleting user
func (srv *UserService) Delete(e context.Context, username string) error {
	return srv.conn.Delete(e, username)
}
