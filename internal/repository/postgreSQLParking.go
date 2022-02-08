package repository

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	log "github.com/sirupsen/logrus"
)

// Add function for inserting a parking lot into sql table
func (rep *PostgresParking) Add(e context.Context, lot *model.ParkingLot) error {
	_, err := rep.PoolParking.Exec(e, "INSERT INTO parking (num,inparking,remark) VALUES ($1,$2,$3)", lot.Num, lot.InParking, lot.Remark)
	if err != nil {
		log.Errorf("can't create parking lot %s", err)
		return err
	}
	return err
}

// GetAll function for getting all parking lots from a sql table
func (rep *PostgresParking) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	rows, err := rep.PoolParking.Query(e, "SELECT * FROM parking")
	if err != nil {
		log.Errorf("can't select all parking lot %s", err)
		return nil, err
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
		lot.Num = int(values[0].(int32))
		lot.InParking = values[1].(bool)
		lot.Remark = values[2].(string)
		lots = append(lots, &lot)
	}
	return lots, err
}

// GetByNum function for getting parking lot by num from a sql table
func (rep *PostgresParking) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	var lot model.ParkingLot
	err := rep.PoolParking.QueryRow(e, "SELECT num,inparking, remark from parking where num=$1", num).Scan(&lot.Num, &lot.InParking, &lot.Remark)
	if err != nil {
		log.Errorf("can't select parking lot %s", err)
		return nil, err
	}
	return &lot, err
}

// Update function for updating parking lot from a sql table
func (rep *PostgresParking) Update(e context.Context, num int, inParking bool, remark string) error {
	_, err := rep.PoolParking.Exec(e, "UPDATE parking SET inparking =$1,remark =$2 WHERE num = $3", inParking, remark, num)
	if err != nil {
		log.Errorf("can't update parking lot %s", err)
		return err
	}
	return err
}

// Delete function for deleting parking lot from a sql table
func (rep *PostgresParking) Delete(e context.Context, num int) error {
	row, err := rep.PoolParking.Exec(e, "DELETE FROM parking where num=$1", num)
	if err != nil {
		log.Errorf("can't delete parking lot %s", err)
		return err
	}
	if row.RowsAffected() != 1 {
		log.Errorf("nothing to delete %s", err)
		return err
	}
	return err
}
