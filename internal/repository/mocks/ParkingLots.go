// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/EgMeln/CRUDentity/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// ParkingLots is an autogenerated mock type for the ParkingLots type
type ParkingLots struct {
	mock.Mock
}

// Add provides a mock function with given fields: e, lot
func (_m *ParkingLots) Add(e context.Context, lot *model.ParkingLot) error {
	ret := _m.Called(e, lot)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.ParkingLot) error); ok {
		r0 = rf(e, lot)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: e, num
func (_m *ParkingLots) Delete(e context.Context, num int) error {
	ret := _m.Called(e, num)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(e, num)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: e
func (_m *ParkingLots) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	ret := _m.Called(e)

	var r0 []*model.ParkingLot
	if rf, ok := ret.Get(0).(func(context.Context) []*model.ParkingLot); ok {
		r0 = rf(e)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ParkingLot)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(e)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNum provides a mock function with given fields: e, num
func (_m *ParkingLots) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	ret := _m.Called(e, num)

	var r0 *model.ParkingLot
	if rf, ok := ret.Get(0).(func(context.Context, int) *model.ParkingLot); ok {
		r0 = rf(e, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ParkingLot)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(e, num)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: e, num, inParking, remark
func (_m *ParkingLots) Update(e context.Context, num int, inParking bool, remark string) error {
	ret := _m.Called(e, num, inParking, remark)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, bool, string) error); ok {
		r0 = rf(e, num, inParking, remark)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
