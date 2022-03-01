package repository

import (
	"context"
	"testing"
	"time"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/gofrs/uuid"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParkingLotCache_GetByNum(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rep := NewParkingLotCache(ctx, redisClient)
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expected := &model.ParkingLot{ID: uuid2.UUID(id), Num: 10, InParking: false, Remark: "remark"}
	err := rep.Add(ctx, expected)
	require.NoError(t, err)
	time.Sleep(time.Second * 1)
	lot, err := rep.GetByNum(ctx, expected.Num)
	require.NoError(t, err)
	require.Equal(t, expected.Num, lot.Num)
	require.Equal(t, expected.InParking, lot.InParking)
	require.Equal(t, expected.Remark, lot.Remark)
}
func TestParkingLotCache_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rep := NewParkingLotCache(ctx, redisClient)
	id, ok := uuid.NewV1()
	require.NoError(t, ok)
	expected := &model.ParkingLot{ID: uuid2.UUID(id), Num: 11, InParking: false, Remark: "remark"}
	err := rep.Add(ctx, expected)
	require.NoError(t, err)
	err = rep.Delete(ctx, expected.Num)
	require.NoError(t, err)
	_, err = rep.GetByNum(ctx, expected.Num)
	require.Error(t, err)
}
