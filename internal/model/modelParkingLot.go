package model

type ParkingLot struct {
	Num       int    `json:"num" form:"num"`
	InParking bool   `json:"in_parking" form:"in_parking"`
	Remark    string `json:"remark" form:"remark"`
}
