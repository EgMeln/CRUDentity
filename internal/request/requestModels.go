package request

type ParkingLotCreate struct {
	Num       int    `json:"num" form:"num"`
	InParking bool   `json:"in_parking" form:"in_parking"`
	Remark    string `json:"remark" form:"remark"`
}

type ParkingLotUpdate struct {
	InParking bool   `json:"in_parking" form:"in_parking"`
	Remark    string `json:"remark" form:"remark"`
}
