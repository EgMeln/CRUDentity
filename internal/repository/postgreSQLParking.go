package repository

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
)

func (rep *Postgres) Add(e context.Context, lot *model.ParkingLot) error {
	_, err := rep.Pool.Exec(e, "INSERT INTO parking (num,inparking,remark) VALUES ($1,$2,$3)", lot.Num, lot.InParking, lot.Remark)
	if err != nil {
		return fmt.Errorf("can't create parking lot %w", err)
	}
	return err
}

func (rep *Postgres) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	rows, err := rep.Pool.Query(e, "SELECT * FROM parking")
	if err != nil {
		return nil, fmt.Errorf("can't select all parking lot %w", err)
	}
	var lots []*model.ParkingLot
	for rows.Next() {
		var lot model.ParkingLot
		values, err := rows.Values()
		if err != nil {
			return lots, err
		}
		lot.Num = int(values[0].(int32))
		lot.InParking = values[1].(bool)
		lot.Remark = values[2].(string)
		lots = append(lots, &lot)
		fmt.Println(lots)
	}
	return lots, err
}

func (rep *Postgres) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	var lot model.ParkingLot
	err := rep.Pool.QueryRow(e, "SELECT num,inparking, remark from parking where num=$1", num).Scan(&lot.Num, &lot.InParking, &lot.Remark)
	if err != nil {
		return nil, fmt.Errorf("can't select parking lot %w", err)
	}
	return &lot, err
}

func (rep *Postgres) Update(e context.Context, num int, inParking bool, remark string) error {
	_, err := rep.Pool.Exec(e, "UPDATE parking SET inparking =$1,remark =$2 WHERE num = $3", inParking, remark, num)
	if err != nil {
		return fmt.Errorf("can't update parking lot %w", err)
	}
	return err
}

func (rep *Postgres) Delete(e context.Context, num int) error {
	row, err := rep.Pool.Exec(e, "DELETE FROM parking where num=$1", num)
	if err != nil {
		return fmt.Errorf("can't delete parking lot %w", err)
	}
	if row.RowsAffected() != 1 {
		return fmt.Errorf("nothing to delete%w", err)
	}
	return err
}
