package repository

import (
	"context"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoUser_Add(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("users")
	rep := MongoUser{collection}
	expected := &model.User{Username: "test", Password: "1", Admin: true}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var user model.User
	err = collection.FindOne(context.Background(), bson.M{"username": expected.Username}).Decode(&user)
	require.NoError(t, err)
	require.Equal(t, expected.Username, user.Username)
	require.Equal(t, expected.Password, user.Password)
	require.Equal(t, expected.Admin, user.Admin)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoUser_GetAll(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("users")
	rep := MongoUser{collection}
	expectedFirst := &model.User{Username: "test", Password: "1", Admin: true}
	expectedSecond := &model.User{Username: "test2", Password: "2", Admin: false}
	_, err := collection.InsertOne(context.Background(), expectedFirst)
	require.NoError(t, err)
	_, err = collection.InsertOne(context.Background(), expectedSecond)
	require.NoError(t, err)
	lots, err := rep.GetAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, expectedFirst.Username, lots[0].Username)
	require.Equal(t, expectedFirst.Password, lots[0].Password)
	require.Equal(t, expectedFirst.Admin, lots[0].Admin)
	require.Equal(t, expectedSecond.Username, lots[1].Username)
	require.Equal(t, expectedSecond.Password, lots[1].Password)
	require.Equal(t, expectedSecond.Admin, lots[1].Admin)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoUser_Get(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("users")
	rep := MongoUser{collection}
	expected := &model.User{Username: "test", Password: "1", Admin: true}
	_, err := collection.InsertOne(context.Background(), expected)
	require.NoError(t, err)
	user, err := rep.Get(context.Background(), expected)
	require.NoError(t, err)
	require.Equal(t, expected.Username, user.Username)
	require.Equal(t, expected.Password, user.Password)
	require.Equal(t, expected.Admin, user.Admin)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoUser_Update(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("users")
	rep := MongoUser{collection}
	expected := &model.User{Username: "test", Password: "1", Admin: true}
	_, err := collection.InsertOne(context.Background(), expected)
	require.NoError(t, err)
	var user model.User
	err = rep.Update(context.Background(), &model.User{Username: expected.Username, Password: "123", Admin: true})
	require.NoError(t, err)
	err = collection.FindOne(context.Background(), bson.M{"username": expected.Username}).Decode(&user)
	require.NoError(t, err)
	require.Equal(t, expected.Username, user.Username)
	require.Equal(t, "123", user.Password)
	require.Equal(t, expected.Admin, user.Admin)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoUser_Delete(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("users")
	rep := MongoUser{collection}
	expected := &model.User{Username: "test", Password: "1", Admin: true}
	_, err := collection.InsertOne(context.Background(), expected)
	require.NoError(t, err)
	err = rep.Delete(context.Background(), expected.Username)
	require.NoError(t, err)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
