package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// ParkingLotCache struct for redis client
type ParkingLotCache struct {
	client      *redis.Client
	parkingMap  map[string]*model.ParkingLot
	redisStream string
}

// ParkingLotCacheRedis struct for cache func
type ParkingLotCacheRedis interface {
	Add(e context.Context, lot *model.ParkingLot) error
	GetByNum(e context.Context, num int) (*model.ParkingLot, error)
	Delete(e context.Context, num int) error
}

// NewParkingLotCache returns new instance of ParkingLotCache
func NewParkingLotCache(ctx context.Context, cln *redis.Client) *ParkingLotCache {
	red := &ParkingLotCache{client: cln, redisStream: "STREAM", parkingMap: make(map[string]*model.ParkingLot)}
	go red.StartProcessing(ctx)
	return red
}

// Add record about cache parking lot
func (red *ParkingLotCache) Add(e context.Context, lot *model.ParkingLot) error {
	err := red.client.XAdd(e, &redis.XAddArgs{
		Stream: red.redisStream,
		Values: map[string]interface{}{
			"num":        lot.Num,
			"in_parking": lot.InParking,
			"remark":     lot.Remark,
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("redis create parking lot error %w", err)
	}
	return nil
}

// GetByNum getting parking lot cache by num
func (red *ParkingLotCache) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	parkingLot, err := red.parkingMap[strconv.Itoa(num)]
	if !err {
		return nil, fmt.Errorf("redis get parking lot cache error %v", err)
	}
	return parkingLot, nil
}

// Delete deleting parking lot cache
func (red *ParkingLotCache) Delete(e context.Context, num int) error {
	delete(red.parkingMap, strconv.Itoa(num))
	return nil
}

// StartProcessing process the received message
func (red *ParkingLotCache) StartProcessing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			streams, err := red.client.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{red.redisStream, "$"},
				Count:   1,
				Block:   0,
			}).Result()
			if err != nil {
				log.Warnf("redis start process error %v", err)
			}
			if streams[0].Messages == nil {
				log.Warn("empty message")
			}
			stream := streams[0].Messages[0].Values
			num, err := strconv.Atoi(stream["num"].(string))
			if err != nil {
				log.Warnf("converting num to int error %v", err)
			}
			inPark, err := strconv.ParseBool(stream["in_parking"].(string))
			if err != nil {
				log.Warnf("converting in_parking to bool error %v", err)
			}
			red.parkingMap[stream["num"].(string)] = &model.ParkingLot{
				Num:       num,
				InParking: inPark,
				Remark:    stream["remark"].(string),
			}
		}
	}
}
