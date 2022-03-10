package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/stretchr/testify/require"
)

func TestPostgresToken_Add(t *testing.T) {
	rep := PostgresToken{postgresDB}
	expected := &model.Token{Username: "test", RefreshToken: "expected"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var token model.Token
	err = postgresDB.QueryRow(context.Background(), "SELECT username,token from tokens where username=$1", expected.Username).Scan(&token.Username, &token.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, expected.Username, token.Username)
	require.Equal(t, expected.RefreshToken, token.RefreshToken)
	_, err = postgresDB.Exec(context.Background(), "truncate tokens")
	require.NoError(t, err)
}
func TestPostgresToken_Get(t *testing.T) {
	rep := PostgresToken{postgresDB}
	expected := &model.Token{Username: "test", RefreshToken: fmt.Sprintf("%x", sha256.Sum256([]byte("expected")))}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO tokens (username,token) VALUES ($1,$2)", expected.Username, expected.RefreshToken)
	require.NoError(t, err)
	token, err := rep.Get(context.Background(), expected.Username)
	require.NoError(t, err)
	require.Equal(t, expected.RefreshToken, token)
	_, err = postgresDB.Exec(context.Background(), "truncate tokens")
	require.NoError(t, err)
}
func TestPostgresToken_Delete(t *testing.T) {
	rep := PostgresToken{postgresDB}
	expected := &model.Token{Username: "test", RefreshToken: fmt.Sprintf("%x", sha256.Sum256([]byte("expected")))}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO tokens (username,token) VALUES ($1,$2)", expected.Username, expected.RefreshToken)
	require.NoError(t, err)
	err = rep.Delete(context.Background(), expected.Username)
	require.NoError(t, err)
	_, err = postgresDB.Exec(context.Background(), "truncate tokens")
	require.NoError(t, err)
}
