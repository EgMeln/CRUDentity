// Package request contain model of request
package request

// ParkingLotCreate struct that contain record info,which will be recorded when creating parking lot
type ParkingLotCreate struct {
	Num       int    `json:"num" form:"num"`
	InParking bool   `json:"in_parking" form:"in_parking"`
	Remark    string `json:"remark" form:"remark"`
}

// ParkingLotUpdate struct that contain record info,which will be recorded when updating parking lot
type ParkingLotUpdate struct {
	InParking bool   `json:"in_parking" form:"in_parking"`
	Remark    string `json:"remark" form:"remark"`
}

// SignInSignUpUser struct that contain record info,which will be recorded when sign up or sign in user
type SignInSignUpUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
}

// UpdateUser struct that contain record info,which will be recorded when updating user
type UpdateUser struct {
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
}

// RefreshToken struct that contain record info,which will be recorded when refreshing token
type RefreshToken struct {
	Username string `json:"username" form:"username"`
	Admin    bool   `json:"admin" form:"admin"`
}
