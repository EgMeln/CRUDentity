package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// ParkingLot struct that contain record info about parking
type ParkingLot struct {
	Num       int    `json:"num" form:"num"`
	InParking bool   `json:"in_parking" form:"in_parking"`
	Remark    string `json:"remark" form:"remark"`
}

const minNumParking = 1000
const maxNumParking = 9999

// Validate checks the validity of the park number
func (park *ParkingLot) Validate() error {
	return validation.ValidateStruct(park,
		validation.Field(&park.Num, validation.Required, validation.Min(minNumParking), validation.Max(maxNumParking)),
		validation.Field(&park.InParking, validation.Required, validation.In("true", "false")),
	)
}
