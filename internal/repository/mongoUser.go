package repository

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add function for inserting a user into mongo table
func (rep *MongoUser) Add(e context.Context, user *model.User) error {
	_, err := rep.CollectionUsers.InsertOne(e, user)
	if err != nil {
		log.Errorf("can't create user %s", err)
		return err
	}
	return err
}

// GetAll function for getting all users from a mongo table
func (rep *MongoUser) GetAll(e context.Context) ([]*model.User, error) { //nolint:dupl //Different business logic
	rows, err := rep.CollectionUsers.Find(e, bson.M{})
	if err != nil {
		log.Errorf("can't select all users %s", err)
		return nil, err
	}
	var users []*model.User
	for rows.Next(e) {
		user := new(model.User)
		if ok := rows.Decode(user); ok != nil {
			return nil, ok
		}
		users = append(users, user)
	}
	if ok := rows.Close(e); ok != nil {
		return nil, ok
	}
	return users, err
}

// Get function for getting user by username from a mongo table
func (rep *MongoUser) Get(e context.Context, username string) (*model.User, error) {
	var user model.User
	err := rep.CollectionUsers.FindOne(e, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		log.Errorf("record doesn't exist %s", err)
		return nil, err
	} else if err != nil {
		log.Errorf("can't select user %s", err)
		return nil, err
	}
	return &user, err
}

// Update function for updating user from a mongo table
func (rep *MongoUser) Update(e context.Context, username, password string, admin bool) error {
	_, err := rep.CollectionUsers.UpdateOne(e, bson.M{"username": username}, bson.M{"$set": bson.M{"password": password, "admin": admin}})
	if err != nil {
		log.Errorf("can't update user %s", err)
		return err
	}
	return err
}

// Delete function for deleting user from a mongo table
func (rep *MongoUser) Delete(e context.Context, username string) error {
	row, err := rep.CollectionUsers.DeleteOne(e, bson.M{"username": username})
	if err != nil {
		log.Errorf("can't delete user %s", err)
		return err
	}
	if row.DeletedCount == 0 {
		log.Errorf("nothing to delete %s", err)
		return err
	}
	return err
}
