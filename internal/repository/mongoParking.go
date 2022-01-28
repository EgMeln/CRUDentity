package repository

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rep *Mongo) Add(e context.Context, lot *model.ParkingLot) error {
	_, err := rep.CollectionParkingLot.InsertOne(e, lot)
	if err != nil {
		return fmt.Errorf("can't create parking lot %w", err)
	}
	return err
}

func (rep *Mongo) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	rows, err := rep.CollectionParkingLot.Find(e, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't select all parking lot %w", err)
	}
	var lots []*model.ParkingLot
	for rows.Next(e) {
		lot := new(model.ParkingLot)
		if err := rows.Decode(lot); err != nil {
			return nil, err
		}
		lots = append(lots, lot)
	}
	if err := rows.Close(e); err != nil {
		return nil, err
	}
	return lots, err
}

func (rep *Mongo) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	var lot model.ParkingLot
	err := rep.CollectionParkingLot.FindOne(e, bson.M{"num": num}).Decode(&lot)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("record doesn't exist %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't select parking lot %w", err)
	}
	return &lot, err
}

func (rep *Mongo) Update(e context.Context, num int, inParking bool, remark string) error {
	_, err := rep.CollectionParkingLot.UpdateOne(e, bson.M{"num": num}, bson.M{"$set": bson.M{"inparking": inParking, "remark": remark}})
	if err != nil {
		return fmt.Errorf("can't update parking lot %w", err)
	}
	return err
}

func (rep *Mongo) Delete(e context.Context, num int) error {
	row, err := rep.CollectionParkingLot.DeleteOne(e, bson.M{"num": num})
	if err != nil {
		return fmt.Errorf("can't delete parking lot %w", err)
	}
	if row.DeletedCount == 0 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
