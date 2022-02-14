// Package request contain model of request
package request

// ParkingLotCreate struct that contain record info,which will be recorded when creating parking lot
type ParkingLotCreate struct {
	Num       int    `json:"num" form:"num" validate:"required"`
	InParking bool   `json:"in_parking" form:"in_parking" validate:"required"`
	Remark    string `json:"remark" form:"remark"`
}

// ParkingLotUpdate struct that contain record info,which will be recorded when updating parking lot
type ParkingLotUpdate struct {
	InParking bool   `json:"in_parking" form:"in_parking" validate:"required"`
	Remark    string `json:"remark" form:"remark"`
}

// SignUpUser struct that contain record info,which will be recorded when sign up user
type SignUpUser struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// SignInUser struct that contain record info,which will be recorded when sign in user
type SignInUser struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// UpdateUser struct that contain record info,which will be recorded when updating user
type UpdateUser struct {
	Password string `json:"password" form:"password" validate:"required"`
}

// RefreshToken struct that contain record info,which will be recorded when refreshing token
type RefreshToken struct {
	Username string `json:"username" form:"username" validate:"required"`
}
