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

type SignInSignUpUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
}

type UpdateUser struct {
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
}
