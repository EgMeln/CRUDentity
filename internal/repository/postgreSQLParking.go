package repository

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
)

// Add function for inserting a parking lot into sql table
func (rep *PostgresParking) Add(e context.Context, lot *model.ParkingLot) error {
	_, err := rep.PoolParking.Exec(e, "INSERT INTO parking (_id,num,inparking,remark) VALUES ($1,$2,$3,$4)", lot.ID, lot.Num, lot.InParking, lot.Remark)
	if err != nil {
		return fmt.Errorf("can't create parking lot %w", err)
	}
	return err
}

// GetAll function for getting all parking lots from a sql table
func (rep *PostgresParking) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	rows, err := rep.PoolParking.Query(e, "SELECT * FROM parking")
	if err != nil {
		return nil, fmt.Errorf("can't select all parking lot %w", err)
	}
	defer rows.Close()
	var lots []*model.ParkingLot
	for rows.Next() {
		var lot model.ParkingLot
		var values []interface{}
		values, err = rows.Values()
		if err != nil {
			return lots, err
		}
		lot.ID = values[0].([16]uint8)
		lot.Num = int(values[1].(int32))
		lot.InParking = values[2].(bool)
		lot.Remark = values[3].(string)
		lots = append(lots, &lot)
	}
	return lots, err
}

// GetByNum function for getting parking lot by num from a sql table
func (rep *PostgresParking) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	var lot model.ParkingLot
	err := rep.PoolParking.QueryRow(e, "SELECT _id,num,inparking, remark from parking where num=$1", num).Scan(&lot.ID, &lot.Num, &lot.InParking, &lot.Remark)
	if err != nil {
		return nil, fmt.Errorf("can't select parking lot %w", err)
	}
	return &lot, err
}

// Update function for updating parking lot from a sql table
func (rep *PostgresParking) Update(e context.Context, num int, inParking bool, remark string) error {
	_, err := rep.PoolParking.Exec(e, "UPDATE parking SET inparking =$1,remark =$2 WHERE num = $3", inParking, remark, num)
	if err != nil {
		return fmt.Errorf("can't update parking lot %w", err)
	}
	return err
}

// Delete function for deleting parking lot from a sql table
func (rep *PostgresParking) Delete(e context.Context, num int) error {
	row, err := rep.PoolParking.Exec(e, "DELETE FROM parking where num=$1", num)
	if err != nil {
		return fmt.Errorf("can't delete parking lot %w", err)
	}
	if row.RowsAffected() != 1 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
