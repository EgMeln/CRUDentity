package model

import "github.com/google/uuid"

// ParkingLot struct that contain record info about parking
type ParkingLot struct {
	ID        uuid.UUID `json:"_id" form:"_id" `
	Num       int       `json:"num" form:"num" `
	InParking bool      `json:"in_parking" form:"in_parking"`
	Remark    string    `json:"remark" form:"remark"`
}
