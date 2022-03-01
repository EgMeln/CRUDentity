package repository

import (
	"context"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/stretchr/testify/require"
)

func TestPostgresUser_Add(t *testing.T) {
	rep := PostgresUser{postgresDB}
	expected := &model.User{Username: "test1", Password: "1", Admin: true}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var user model.User
	row := postgresDB.QueryRow(context.Background(), "SELECT username,password,is_admin from users where username=$1", expected.Username)
	err = row.Scan(&user.Username, &user.Password, &user.Admin)
	require.NoError(t, err)
	require.Equal(t, expected.Username, user.Username)
	require.Equal(t, expected.Password, user.Password)
	require.Equal(t, expected.Admin, user.Admin)
	_, err = postgresDB.Exec(context.Background(), "truncate users")
	require.NoError(t, err)
}
func TestPostgresUser_GetAll(t *testing.T) {
	rep := PostgresUser{postgresDB}
	expectedFirst := &model.User{Username: "test2", Password: "1", Admin: true}
	expectedSecond := &model.User{Username: "test3", Password: "2", Admin: false}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO users (username,password,is_admin) VALUES ($1,$2,$3)",
		expectedFirst.Username, expectedFirst.Password, expectedFirst.Admin)
	require.NoError(t, err)
	_, err = postgresDB.Exec(context.Background(), "INSERT INTO users (username,password,is_admin) VALUES ($1,$2,$3)",
		expectedSecond.Username, expectedSecond.Password, expectedSecond.Admin)
	require.NoError(t, err)
	lots, ok := rep.GetAll(context.Background())
	require.NoError(t, ok)
	require.Equal(t, expectedFirst.Username, lots[0].Username)
	require.Equal(t, expectedFirst.Password, lots[0].Password)
	require.Equal(t, expectedFirst.Admin, lots[0].Admin)
	require.Equal(t, expectedSecond.Username, lots[1].Username)
	require.Equal(t, expectedSecond.Password, lots[1].Password)
	require.Equal(t, expectedSecond.Admin, lots[1].Admin)
	_, err = postgresDB.Exec(context.Background(), "truncate users")
	require.NoError(t, err)
}
func TestPostgresUser_Get(t *testing.T) {
	rep := PostgresUser{postgresDB}
	expected := &model.User{Username: "test4", Password: "1", Admin: true}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO users (username,password,is_admin) VALUES ($1,$2,$3)", expected.Username, expected.Password, expected.Admin)
	require.NoError(t, err)
	user, ok := rep.Get(context.Background(), expected)
	require.NoError(t, ok)
	require.Equal(t, expected.Username, user.Username)
	require.Equal(t, expected.Password, user.Password)
	require.Equal(t, expected.Admin, user.Admin)
	_, err = postgresDB.Exec(context.Background(), "truncate users")
	require.NoError(t, err)
}
func TestPostgresUser_Update(t *testing.T) {
	rep := PostgresUser{postgresDB}
	expected := &model.User{Username: "test5", Password: "1", Admin: true}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO users (username,password,is_admin) VALUES ($1,$2,$3)", expected.Username, expected.Password, expected.Admin)
	require.NoError(t, err)
	err = rep.Update(context.Background(), expected)
	require.NoError(t, err)
	var user model.User
	row := postgresDB.QueryRow(context.Background(), "SELECT username,password,is_admin from users where username=$1", expected.Username)
	err = row.Scan(&user.Username, &user.Password, &user.Admin)
	require.NoError(t, err)
	require.Equal(t, expected.Username, user.Username)
	require.Equal(t, expected.Password, user.Password)
	require.Equal(t, expected.Admin, user.Admin)
	_, err = postgresDB.Exec(context.Background(), "truncate users")
	require.NoError(t, err)
}
func TestPostgresUser_Delete(t *testing.T) {
	rep := PostgresUser{postgresDB}
	expected := &model.User{Username: "test6", Password: "1", Admin: true}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO users (username,password,is_admin) VALUES ($1,$2,$3)", expected.Username, expected.Password, expected.Admin)
	require.NoError(t, err)
	err = rep.Delete(context.Background(), expected.Username)
	require.NoError(t, err)
	_, err = postgresDB.Exec(context.Background(), "truncate users")
	require.NoError(t, err)
}
