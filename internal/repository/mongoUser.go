package repository

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add function for inserting a user into mongo table
func (rep *MongoUser) Add(e context.Context, user *model.User) error {
	_, err := rep.CollectionUsers.InsertOne(e, user)
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}

// GetAll function for getting all users from a mongo table
func (rep *MongoUser) GetAll(e context.Context) ([]*model.User, error) {
	rows, err := rep.CollectionUsers.Find(e, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't select all users %w", err)
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
		return nil, fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't select user %w", err)
	}
	return &user, err
}

// Update function for updating user from a mongo table
func (rep *MongoUser) Update(e context.Context, username, password string, admin bool) error {
	_, err := rep.CollectionUsers.UpdateOne(e, bson.M{"username": username}, bson.M{"$set": bson.M{"password": password, "admin": admin}})
	if err != nil {
		return fmt.Errorf("can't update user %w", err)
	}
	return err
}

// Delete function for deleting user from a mongo table
func (rep *MongoUser) Delete(e context.Context, username string) error {
	row, err := rep.CollectionUsers.DeleteOne(e, bson.M{"username": username})
	if err != nil {
		return fmt.Errorf("can't delete user %w", err)
	}
	if row.DeletedCount == 0 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
