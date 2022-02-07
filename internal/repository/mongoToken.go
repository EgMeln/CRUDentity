package repository

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rep *MongoToken) Add(e context.Context, token *model.Token) error {
	_, err := rep.CollectionTokens.InsertOne(e, token)
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}
func (rep *MongoToken) Get(e context.Context, username string) (string, error) {
	var token model.Token
	err := rep.CollectionTokens.FindOne(e, bson.M{"username": username}).Decode(&token)
	if err == mongo.ErrNoDocuments {
		return "", fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return "", fmt.Errorf("can't select token %w", err)
	}
	return token.RefreshToken, err
}
func (rep *MongoToken) Delete(e context.Context, username string) error {
	row, err := rep.CollectionTokens.DeleteOne(e, bson.M{"username": username})
	if err != nil {
		return fmt.Errorf("can't delete user %w", err)
	}
	if row.DeletedCount == 0 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
