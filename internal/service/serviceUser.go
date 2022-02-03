package service

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

type UserService struct {
	conn repository.Users
}

func NewUserServicePostgres(rep *repository.Postgres) *UserService {
	return &UserService{conn: rep}
}
func NewUserServiceMongo(rep *repository.Mongo) *UserService {
	return &UserService{conn: rep}
}
func (srv *UserService) Add(e context.Context, user *model.User) error {
	return srv.conn.AddUser(e, user)
}
func (srv *UserService) GetAll(e context.Context) ([]*model.User, error) {
	return srv.conn.GetAllUser(e)
}
func (srv *UserService) Get(e context.Context, username string) (*model.User, error) {
	return srv.conn.GetUser(e, username)
}
func (srv *UserService) Update(e context.Context, username string, password string, admin bool) error {
	return srv.conn.UpdateUser(e, username, password, admin)
}
func (srv *UserService) Delete(e context.Context, username string) error {
	return srv.conn.DeleteUser(e, username)
}
func (srv *UserService) AddToken(e context.Context, username string, tokenStr string) error {
	return srv.conn.AddToken(e, username, tokenStr)
}
func (srv *UserService) GetToken(e context.Context, username string) (string, error) {
	return srv.conn.GetToken(e, username)
}
