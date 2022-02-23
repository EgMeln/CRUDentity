package repository

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// ParkingLotCache struct for redis client
type ParkingLotCache struct {
	client      *redis.Client
	parkingMap  map[int]*model.ParkingLot
	redisStream string
	mu          sync.RWMutex
}

// ParkingLotCacheRedis struct for cache func
type ParkingLotCacheRedis interface {
	Add(e context.Context, lot *model.ParkingLot) error
	GetByNum(e context.Context, num int) (*model.ParkingLot, error)
	Delete(e context.Context, num int) error
}

// NewParkingLotCache returns new instance of ParkingLotCache
func NewParkingLotCache(ctx context.Context, cln *redis.Client) *ParkingLotCache {
	red := &ParkingLotCache{client: cln, redisStream: "STREAM", parkingMap: make(map[int]*model.ParkingLot), mu: sync.RWMutex{}}
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
	red.mu.Lock()
	parkingLot, err := red.parkingMap[num]
	red.mu.Unlock()
	if !err {
		return nil, fmt.Errorf("redis get parking lot cache error %v", err)
	}
	return parkingLot, nil
}

// Delete deleting parking lot cache
func (red *ParkingLotCache) Delete(e context.Context, num int) error {
	red.mu.Lock()
	res := red.client.XDel(e, red.redisStream, strconv.Itoa(num))
	if res == nil {
		return fmt.Errorf("nothing to delete from redis stream %v", res)
	}
	delete(red.parkingMap, num)
	red.mu.Unlock()
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
			red.mu.Lock()
			red.parkingMap[num] = &model.ParkingLot{
				Num:       num,
				InParking: inPark,
				Remark:    stream["remark"].(string),
			}
			red.mu.Unlock()
		}
	}
}
