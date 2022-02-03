package repository

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rep *Mongo) AddUser(e context.Context, user *model.User) error {
	_, err := rep.CollectionUsers.InsertOne(e, user)
	if err != nil {
		return fmt.Errorf("can't create user %w", err)
	}
	return err
}
func (rep *Mongo) GetAllUser(e context.Context) ([]*model.User, error) {
	rows, err := rep.CollectionUsers.Find(e, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't select all users %w", err)
	}
	defer rows.Close(e)
	var users []*model.User
	for rows.Next(e) {
		lot := new(model.User)
		if err := rows.Decode(lot); err != nil {
			return nil, err
		}
		users = append(users, lot)
	}
	if err := rows.Close(e); err != nil {
		return nil, err
	}
	return users, err
}
func (rep *Mongo) GetUser(e context.Context, username string) (*model.User, error) {
	var user model.User
	err := rep.CollectionUsers.FindOne(e, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't select user %w", err)
	}
	return &user, err
}
func (rep *Mongo) UpdateUser(e context.Context, username string, password string, admin bool) error {
	_, err := rep.CollectionUsers.UpdateOne(e, bson.M{"username": username}, bson.M{"$set": bson.M{"password": password, "admin": admin}})
	if err != nil {
		return fmt.Errorf("can't update user %w", err)
	}
	return err
}
func (rep *Mongo) DeleteUser(e context.Context, username string) error {
	row, err := rep.CollectionUsers.DeleteOne(e, bson.M{"username": username})
	if err != nil {
		return fmt.Errorf("can't delete user %w", err)
	}
	if row.DeletedCount == 0 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
