package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoToken_Add(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("tokens")
	rep := MongoToken{collection}
	expected := &model.Token{Username: "test", RefreshToken: "expected"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var token model.Token
	err = collection.FindOne(context.Background(), bson.M{"username": expected.Username}).Decode(&token)
	require.NoError(t, err)
	require.Equal(t, expected.Username, token.Username)
	require.Equal(t, expected.RefreshToken, token.RefreshToken)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoToken_Get(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("tokens")
	rep := MongoToken{collection}
	expected := &model.Token{Username: "test", RefreshToken: fmt.Sprintf("%x", sha256.Sum256([]byte("expected")))}
	_, err := collection.InsertOne(context.Background(), expected)
	require.NoError(t, err)
	token, err := rep.Get(context.Background(), expected.Username)
	require.NoError(t, err)
	require.Equal(t, expected.RefreshToken, token)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoToken_Delete(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("tokens")
	rep := MongoToken{collection}
	expected := &model.Token{Username: "test", RefreshToken: fmt.Sprintf("%x", sha256.Sum256([]byte("expected")))}
	_, err := collection.InsertOne(context.Background(), expected)
	require.NoError(t, err)
	err = rep.Delete(context.Background(), expected.Username)
	require.NoError(t, err)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
