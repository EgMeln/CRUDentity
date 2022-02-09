package repository

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add function for inserting a token into sql table
func (rep *PostgresToken) Add(e context.Context, token *model.Token) error {
	_, err := rep.PoolToken.Exec(e, "INSERT INTO tokens (username,admin,token) VALUES ($1,$2,$3)", token.Username, token.Admin, token.RefreshToken)

	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}

// Get function for getting token from a sql table
func (rep *PostgresToken) Get(e context.Context, username string) (string, error) {
	var token model.Token
	err := rep.PoolToken.QueryRow(e, "SELECT username,admin,token from tokens where username=$1", username).Scan(&token.Username, &token.Admin, &token.RefreshToken)
	if err == mongo.ErrNoDocuments {
		return "", fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return "", fmt.Errorf("can't select token %w", err)
	}
	return token.RefreshToken, err
}

// Delete function for deleting token from a sql table
func (rep *PostgresToken) Delete(e context.Context, username string) error {
	row, err := rep.PoolToken.Exec(e, "DELETE FROM tokens where username=$1", username)
	if err != nil {
		return fmt.Errorf("can't delete token %w", err)
	}
	if row.RowsAffected() != 1 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
