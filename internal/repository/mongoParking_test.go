package repository

import (
	"context"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoParking_Add(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("egormelnikov")
	rep := MongoParking{collection}
	expected := &model.ParkingLot{Num: 1, InParking: false, Remark: "remark"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var lot model.ParkingLot
	err = rep.CollectionParkingLot.FindOne(context.Background(), bson.M{"num": expected.Num}).Decode(&lot)
	require.NoError(t, err)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, expected.InParking, lot.InParking)
	require.Equal(t, expected.Remark, lot.Remark)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoParking_GetAll(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("egormelnikov")
	rep := MongoParking{collection}
	expectedFirst := &model.ParkingLot{Num: 1, InParking: false, Remark: "remark1"}
	expectedSecond := &model.ParkingLot{Num: 2, InParking: true, Remark: "remark2"}
	err := rep.Add(context.Background(), expectedFirst)
	require.NoError(t, err)
	err = rep.Add(context.Background(), expectedSecond)
	require.NoError(t, err)
	lots, err := rep.GetAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, expectedFirst.Num, lots[0].Num)
	require.Equal(t, expectedFirst.InParking, lots[0].InParking)
	require.Equal(t, expectedFirst.Remark, lots[0].Remark)
	require.Equal(t, expectedSecond.Num, lots[1].Num)
	require.Equal(t, expectedSecond.InParking, lots[1].InParking)
	require.Equal(t, expectedSecond.Remark, lots[1].Remark)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoParking_GetByNum(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("egormelnikov")
	rep := MongoParking{collection}
	expected := &model.ParkingLot{Num: 1, InParking: false, Remark: "remark"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	lot, err := rep.GetByNum(context.Background(), expected.Num)
	require.NoError(t, err)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, expected.InParking, lot.InParking)
	require.Equal(t, expected.Remark, lot.Remark)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoParking_Update(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("egormelnikov")
	rep := MongoParking{collection}
	expected := &model.ParkingLot{Num: 1, InParking: false, Remark: "remark"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var lot model.ParkingLot
	err = rep.Update(context.Background(), expected.Num, true, "f")
	require.NoError(t, err)
	err = collection.FindOne(context.Background(), bson.M{"num": expected.Num}).Decode(&lot)
	require.NoError(t, err)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, true, lot.InParking)
	require.Equal(t, "f", lot.Remark)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
func TestMongoParking_Delete(t *testing.T) {
	collection := dbClient.Database("egormelnikovdb").Collection("egormelnikov")
	rep := MongoParking{collection}
	expected := &model.ParkingLot{Num: 1, InParking: false, Remark: "remark"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	err = rep.Delete(context.Background(), expected.Num)
	require.NoError(t, err)
	err = collection.Drop(context.Background())
	require.NoError(t, err)
}
