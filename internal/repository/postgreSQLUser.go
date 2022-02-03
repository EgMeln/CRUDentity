package repository

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
)

func (rep *Postgres) AddUser(e context.Context, user *model.User) error {
	_, err := rep.Pool.Exec(e, "INSERT INTO users (username,password,admin) VALUES ($1,$2,$3)", user.Username, user.Password, user.Admin)
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}
func (rep *Postgres) GetAllUser(e context.Context) ([]*model.User, error) {
	rows, err := rep.Pool.Query(e, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("can't select all users %w", err)
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		values, err := rows.Values()
		if err != nil {
			return users, err
		}
		user.Username = values[0].(string)
		user.Password = values[1].(string)
		user.Admin = values[2].(bool)
		users = append(users, &user)
	}

	return users, err
}
func (rep *Postgres) GetUser(e context.Context, username string) (*model.User, error) {
	var user model.User
	err := rep.Pool.QueryRow(e, "SELECT username,password,admin from users where username=$1", username).Scan(&user.Username, &user.Password, &user.Admin)
	if err != nil {
		return nil, fmt.Errorf("can't select parking lot %w", err)
	}
	return &user, err
}
func (rep *Postgres) UpdateUser(e context.Context, username string, password string, admin bool) error {
	_, err := rep.Pool.Exec(e, "UPDATE users SET password =$1,admin =$2 WHERE username = $3", password, admin, username)
	if err != nil {
		return fmt.Errorf("can't update parking lot %w", err)
	}
	return err
}
func (rep *Postgres) DeleteUser(e context.Context, username string) error {
	row, err := rep.Pool.Exec(e, "DELETE FROM users where username=$1", username)
	if err != nil {
		return fmt.Errorf("can't delete parking lot %w", err)
	}
	if row.RowsAffected() != 1 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
