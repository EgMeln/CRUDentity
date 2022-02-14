package repository

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
)

// Add function for inserting a user into sql table
func (rep *PostgresUser) Add(e context.Context, user *model.User) error {
	_, err := rep.PoolUser.Exec(e, "INSERT INTO users (username,password,admin) VALUES ($1,$2,$3)", user.Username, user.Password, user.Admin)
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}

// GetAll function for getting all users from a sql table
func (rep *PostgresUser) GetAll(e context.Context) ([]*model.User, error) {
	rows, err := rep.PoolUser.Query(e, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("can't select all users %w", err)
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		var values []interface{}
		values, err = rows.Values()
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

// Get function for getting user by username from a sql table
func (rep *PostgresUser) Get(e context.Context, username string) (*model.User, error) {
	var user model.User
	err := rep.PoolUser.QueryRow(e, "SELECT username,password,admin from users where username=$1", username).Scan(&user.Username, &user.Password, &user.Admin)
	if err != nil {
		return nil, fmt.Errorf("can't select parking lot %w", err)
	}
	return &user, err
}

// Update function for updating user from a sql table
func (rep *PostgresUser) Update(e context.Context, user *model.User) error {
	_, err := rep.PoolUser.Exec(e, "UPDATE users SET password =$1,admin =$2 WHERE username = $3", user.Password, user.Admin, user.Username)
	if err != nil {
		return fmt.Errorf("can't update parking lot %w", err)
	}
	return err
}

// Delete function for deleting user from a sql table
func (rep *PostgresUser) Delete(e context.Context, username string) error {
	row, err := rep.PoolUser.Exec(e, "DELETE FROM users where username=$1", username)
	if err != nil {
		return fmt.Errorf("can't delete parking lot %w", err)
	}
	if row.RowsAffected() != 1 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
