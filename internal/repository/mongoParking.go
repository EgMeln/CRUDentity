package repository

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add function for inserting a parking lot into mongo table
func (rep *MongoParking) Add(e context.Context, lot *model.ParkingLot) error {
	_, err := rep.CollectionParkingLot.InsertOne(e, lot)
	if err != nil {
		log.Errorf("can't create parking lot %s", err)
		return err
	}
	return err
}

// GetAll function for getting all parking lots from a mongo table
func (rep *MongoParking) GetAll(e context.Context) ([]*model.ParkingLot, error) { //nolint:dupl //Different business logic
	rows, err := rep.CollectionParkingLot.Find(e, bson.M{})
	if err != nil {
		log.Errorf("can't select all parking lots %s", err)
		return nil, err
	}
	var lots []*model.ParkingLot
	for rows.Next(e) {
		lot := new(model.ParkingLot)
		if ok := rows.Decode(lot); ok != nil {
			return nil, ok
		}
		lots = append(lots, lot)
	}
	if ok := rows.Close(e); ok != nil {
		return nil, ok
	}
	return lots, err
}

// GetByNum function for getting parking lot by num from a mongo table
func (rep *MongoParking) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	var lot model.ParkingLot
	err := rep.CollectionParkingLot.FindOne(e, bson.M{"num": num}).Decode(&lot)
	if err == mongo.ErrNoDocuments {
		log.Errorf("record doesn't exist %s", err)
		return nil, err
	} else if err != nil {
		log.Errorf("can't select parking lot %s", err)
		return nil, err
	}
	return &lot, err
}

// Update function for updating parking lot from a mongo table
func (rep *MongoParking) Update(e context.Context, num int, inParking bool, remark string) error {
	_, err := rep.CollectionParkingLot.UpdateOne(e, bson.M{"num": num}, bson.M{"$set": bson.M{"inparking": inParking, "remark": remark}})
	if err != nil {
		log.Errorf("can't update parking lot %s", err)
		return err
	}
	return err
}

// Delete function for deleting parking lot from a mongo table
func (rep *MongoParking) Delete(e context.Context, num int) error {
	row, err := rep.CollectionParkingLot.DeleteOne(e, bson.M{"num": num})
	if err != nil {
		log.Errorf("can't delete parking lot %s", err)
		return err
	}
	if row.DeletedCount == 0 {
		log.Errorf("nothing to delete%s", err)
		return err
	}
	return err
}
