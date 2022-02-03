package service

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

type UserService struct {
	conn repository.Users
}

func NewUserServicePostgres(rep *repository.PostgresUser) *UserService {
	return &UserService{conn: rep}
}
func NewUserServiceMongo(rep *repository.MongoUser) *UserService {
	return &UserService{conn: rep}
}
func (srv *UserService) Add(e context.Context, user *model.User) error {
	return srv.conn.Add(e, user)
}
func (srv *UserService) GetAll(e context.Context) ([]*model.User, error) {
	return srv.conn.GetAll(e)
}
func (srv *UserService) Get(e context.Context, username string) (*model.User, error) {
	return srv.conn.Get(e, username)
}
func (srv *UserService) Update(e context.Context, username string, password string, admin bool) error {
	return srv.conn.Update(e, username, password, admin)
}
func (srv *UserService) Delete(e context.Context, username string) error {
	return srv.conn.Delete(e, username)
}
