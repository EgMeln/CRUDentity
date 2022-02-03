package repository

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rep *Mongo) GetToken(e context.Context, username string) (string, error) {
	var user model.User
	err := rep.CollectionUsers.FindOne(e, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return "", fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return "", fmt.Errorf("can't select token %w", err)
	}
	return user.Token, err
}

func (rep *Mongo) AddToken(e context.Context, username string, tokenStr string) error {
	_, err := rep.CollectionUsers.UpdateOne(e, bson.M{"username": username}, bson.M{"$set": bson.M{"refreshToken": tokenStr}})
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}
