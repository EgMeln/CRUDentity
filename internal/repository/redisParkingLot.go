package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/go-redis/redis/v8"
)

// ParkingLotCache struct for redis client
type ParkingLotCache struct {
	client *redis.Client
}

// ParkingLotCacheRedis struct for cacge func
type ParkingLotCacheRedis interface {
	Add(e context.Context, lot *model.ParkingLot) error
	GetByNum(e context.Context, num int) (*model.ParkingLot, error)
	Delete(e context.Context, num int) error
}

// NewParkingLotCache returns new instance of ParkingLotCache
func NewParkingLotCache(cln *redis.Client) *ParkingLotCache {
	return &ParkingLotCache{client: cln}
}

// Add record about cache parking lot
func (red *ParkingLotCache) Add(e context.Context, lot *model.ParkingLot) error {
	parkingLot, err := json.Marshal(lot)
	if err != nil {
		return fmt.Errorf("json marshal error %w", err)
	}
	if err := red.client.Set(e, strconv.Itoa(lot.Num), parkingLot, time.Minute).Err(); err != nil {
		return fmt.Errorf("redis create parking lot error %w", err)
	}
	return nil
}

// GetByNum getting parking lot cache by num
func (red *ParkingLotCache) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	res, err := red.client.Get(e, strconv.Itoa(num)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get parking lot cache error %w", err)
	}
	var parkingLot model.ParkingLot
	err = json.Unmarshal([]byte(res), &parkingLot)
	if err != nil {
		return nil, fmt.Errorf("redis get parking lot json unmarshal error %w", err)
	}
	return &parkingLot, err
}

// Delete deleting parking lot cache
func (red *ParkingLotCache) Delete(e context.Context, num int) error {
	res, err := red.client.Get(e, strconv.Itoa(num)).Result()
	if err != nil {
		return fmt.Errorf("redis get parking lot cache error %w", err)
	}
	var parkingLot model.ParkingLot
	err = json.Unmarshal([]byte(res), &parkingLot)
	if err != nil {
		return fmt.Errorf("redis get parking lot json unmarshal error %w", err)
	}
	if err := red.client.Del(e, strconv.Itoa(num)).Err(); err != nil {
		return fmt.Errorf("redis delete human cache info error %w", err)
	}
	return nil
}
