package repository

import (
	"context"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/gofrs/uuid"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPostgresParking_Add(t *testing.T) {
	rep := PostgresParking{postgresDB}
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expected := &model.ParkingLot{ID: uuid2.UUID(id), Num: 1, InParking: false, Remark: "remark"}
	err := rep.Add(context.Background(), expected)
	require.NoError(t, err)
	var lot model.ParkingLot
	err = postgresDB.QueryRow(context.Background(), "SELECT _id,num,inparking, remark from parking where num=$1", expected.Num).Scan(&lot.ID, &lot.Num, &lot.InParking, &lot.Remark)
	require.NoError(t, err)
	require.Equal(t, expected.ID, lot.ID)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, expected.InParking, lot.InParking)
	require.Equal(t, expected.Remark, lot.Remark)
	_, err = postgresDB.Exec(context.Background(), "truncate parking")
	require.NoError(t, err)
}
func TestPostgresParking_GetAll(t *testing.T) {
	rep := PostgresParking{postgresDB}
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expectedFirst := &model.ParkingLot{ID: uuid2.UUID(id), Num: 2, InParking: false, Remark: "remark1"}
	id, ok = uuid.NewV1()
	require.NoError(t, ok)
	expectedSecond := &model.ParkingLot{ID: uuid2.UUID(id), Num: 3, InParking: true, Remark: "remark2"}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO parking (_id,num,inparking,remark) VALUES ($1,$2,$3,$4)",
		expectedFirst.ID, expectedFirst.Num, expectedFirst.InParking, expectedFirst.Remark)
	require.NoError(t, err)
	_, err = postgresDB.Exec(context.Background(), "INSERT INTO parking (_id,num,inparking,remark) VALUES ($1,$2,$3,$4)",
		expectedSecond.ID, expectedSecond.Num, expectedSecond.InParking, expectedSecond.Remark)
	require.NoError(t, err)
	var lots []*model.ParkingLot
	lots, err = rep.GetAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, expectedFirst.ID, lots[0].ID)
	require.Equal(t, expectedFirst.Num, lots[0].Num)
	require.Equal(t, expectedFirst.InParking, lots[0].InParking)
	require.Equal(t, expectedFirst.Remark, lots[0].Remark)
	require.Equal(t, expectedSecond.ID, lots[1].ID)
	require.Equal(t, expectedSecond.Num, lots[1].Num)
	require.Equal(t, expectedSecond.InParking, lots[1].InParking)
	require.Equal(t, expectedSecond.Remark, lots[1].Remark)
	_, err = postgresDB.Exec(context.Background(), "truncate parking")
	require.NoError(t, err)
}
func TestPostgresParking_GetByNum(t *testing.T) {
	rep := PostgresParking{postgresDB}
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expected := &model.ParkingLot{ID: uuid2.UUID(id), Num: 4, InParking: false, Remark: "remark"}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO parking (_id,num,inparking,remark) VALUES ($1,$2,$3,$4)",
		expected.ID, expected.Num, expected.InParking, expected.Remark)
	require.NoError(t, err)
	var lot *model.ParkingLot
	lot, err = rep.GetByNum(context.Background(), expected.Num)
	require.NoError(t, err)
	require.Equal(t, expected.ID, lot.ID)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, expected.InParking, lot.InParking)
	require.Equal(t, expected.Remark, lot.Remark)
	_, err = postgresDB.Exec(context.Background(), "truncate parking")
	require.NoError(t, err)
}
func TestPostgresParking_Update(t *testing.T) {
	rep := PostgresParking{postgresDB}
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expected := &model.ParkingLot{ID: uuid2.UUID(id), Num: 5, InParking: false, Remark: "remark"}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO parking (_id,num,inparking,remark) VALUES ($1,$2,$3,$4)", expected.ID, expected.Num, true, "no remark")
	require.NoError(t, err)
	err = rep.Update(context.Background(), expected.Num, expected.InParking, expected.Remark)
	require.NoError(t, err)
	var lot model.ParkingLot
	err = postgresDB.QueryRow(context.Background(), "SELECT _id,num,inparking, remark from parking where num=$1", expected.Num).Scan(&lot.ID, &lot.Num, &lot.InParking, &lot.Remark)
	require.NoError(t, err)
	require.Equal(t, expected.ID, lot.ID)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, expected.InParking, lot.InParking)
	require.Equal(t, expected.Remark, lot.Remark)
	_, err = postgresDB.Exec(context.Background(), "truncate parking")
	require.NoError(t, err)
}
func TestPostgresParking_Delete(t *testing.T) {
	rep := PostgresParking{postgresDB}
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expected := &model.ParkingLot{ID: uuid2.UUID(id), Num: 6, InParking: false, Remark: "remark"}
	_, err := postgresDB.Exec(context.Background(), "INSERT INTO parking (_id,num,inparking,remark) VALUES ($1,$2,$3,$4)",
		expected.ID, expected.Num, expected.InParking, expected.Remark)
	require.NoError(t, err)
	err = rep.Delete(context.Background(), expected.Num)
	require.NoError(t, err)
	_, err = postgresDB.Exec(context.Background(), "truncate parking")
	require.NoError(t, err)
}
