package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rep *Postgres) GetToken(e context.Context, username string) (string, error) {
	var token string
	err := rep.Pool.QueryRow(e, "SELECT token from users where username=$1", username).Scan(&token)
	if err == mongo.ErrNoDocuments {
		return "", fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return "", fmt.Errorf("can't select token %w", err)
	}
	return token, err
}

func (rep *Postgres) AddToken(e context.Context, username string, tokenStr string) error {
	_, err := rep.Pool.Exec(e, "UPDATE users SET token =$1, WHERE username = $3", tokenStr, username)
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}
