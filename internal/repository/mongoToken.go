package repository

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add function for inserting a token into mongo table
func (rep *MongoToken) Add(e context.Context, token *model.Token) error {
	_, err := rep.CollectionTokens.InsertOne(e, token)
	if err != nil {
		log.Errorf("can't create user %s", err)
		return err
	}
	return err
}

// Get function for getting item from a mongo table
func (rep *MongoToken) Get(e context.Context, username string) (string, error) {
	var token model.Token
	err := rep.CollectionTokens.FindOne(e, bson.M{"username": username}).Decode(&token)
	if err == mongo.ErrNoDocuments {
		log.Errorf("record doesn't exist %s", err)
		return "", err
	} else if err != nil {
		log.Errorf("can't select token %s", err)
		return "", err
	}
	return token.RefreshToken, err
}

// Delete function for deleting token from a mongo table
func (rep *MongoToken) Delete(e context.Context, username string) error {
	row, err := rep.CollectionTokens.DeleteOne(e, bson.M{"username": username})
	if err != nil {
		log.Errorf("can't delete user %s", err)
		return err
	}
	if row.DeletedCount == 0 {
		log.Errorf("nothing to delete%s", err)
		return err
	}
	return err
}
