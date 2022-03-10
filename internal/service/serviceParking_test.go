package service

import (
	"context"
	"testing"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MyMockedParking struct {
	parkingLot      mocks.ParkingLots
	parkingLotCache mocks.ParkingLotCacheRedis
}

func TestParkingService_Add(t *testing.T) {
	testParking := new(MyMockedParking)
	id, _ := uuid.NewUUID()
	parkingLot := &model.ParkingLot{ID: id, Num: 1111, InParking: true, Remark: "rem"}
	testParking.parkingLot.On("Add", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.ParkingLot")).Return(
		func(e context.Context, lot *model.ParkingLot) error {
			return nil
		})
	testParking.parkingLotCache.On("Add", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.ParkingLot")).Return(
		func(e context.Context, lot *model.ParkingLot) error {
			return nil
		})
	serviceParking := NewParkingLotService(&testParking.parkingLot, &testParking.parkingLotCache)
	err := serviceParking.Add(context.Background(), parkingLot)
	require.NoError(t, err)
}
func TestParkingService_GetAll(t *testing.T) {
	testParking := new(MyMockedParking)
	testParking.parkingLot.On("GetAll", mock.AnythingOfType("*context.emptyCtx")).Return(
		func(e context.Context) []*model.ParkingLot {
			return []*model.ParkingLot{}
		},
		func(e context.Context) error {
			return nil
		})
	serviceParking := NewParkingLotService(&testParking.parkingLot, &testParking.parkingLotCache)
	_, err := serviceParking.GetAll(context.Background())
	require.NoError(t, err)
}
func TestParkingService_GetByNum(t *testing.T) {
	testParking := new(MyMockedParking)
	id, _ := uuid.NewUUID()
	parkingLot := &model.ParkingLot{ID: id, Num: 1111, InParking: true, Remark: "rem"}
	testParking.parkingLotCache.On("GetByNum", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int")).Return(
		func(e context.Context, num int) *model.ParkingLot {
			return &model.ParkingLot{ID: id, Num: 1111, InParking: true, Remark: "rem"}
		},
		func(e context.Context, num int) error {
			return nil
		})
	testParking.parkingLot.On("GetByNum", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int")).Return(
		func(e context.Context, num int) *model.ParkingLot {
			return &model.ParkingLot{ID: id, Num: 1111, InParking: true, Remark: "rem"}
		},
		func(e context.Context, num int) error {
			return nil
		})
	serviceParking := NewParkingLotService(&testParking.parkingLot, &testParking.parkingLotCache)
	_, err := serviceParking.GetByNum(context.Background(), parkingLot.Num)
	require.NoError(t, err)
}
func TestParkingService_Delete(t *testing.T) {
	testParking := new(MyMockedParking)
	id, _ := uuid.NewUUID()
	parkingLot := &model.ParkingLot{ID: id, Num: 1111, InParking: true, Remark: "rem"}
	testParking.parkingLot.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int")).Return(
		func(e context.Context, num int) error {
			return nil
		})
	testParking.parkingLotCache.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int")).Return(
		func(e context.Context, num int) error {
			return nil
		})
	serviceParking := NewParkingLotService(&testParking.parkingLot, &testParking.parkingLotCache)
	err := serviceParking.Delete(context.Background(), parkingLot.Num)
	require.NoError(t, err)
}
func TestParkingService_Update(t *testing.T) {
	testParking := new(MyMockedParking)
	id, _ := uuid.NewUUID()
	parkingLot := &model.ParkingLot{ID: id, Num: 1111, InParking: true, Remark: "rem"}
	testParking.parkingLot.On("Update", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int"), mock.AnythingOfType("bool"), mock.AnythingOfType("string")).Return(
		func(e context.Context, num int, inParking bool, remark string) error {
			return nil
		})
	testParking.parkingLotCache.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int")).Return(
		func(e context.Context, num int) error {
			return nil
		})
	serviceParking := NewParkingLotService(&testParking.parkingLot, &testParking.parkingLotCache)
	err := serviceParking.Update(context.Background(), parkingLot.Num, parkingLot.InParking, parkingLot.Remark)
	require.NoError(t, err)
}
