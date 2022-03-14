package server

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/EgMeln/CRUDentity/protocol"
	log "github.com/sirupsen/logrus"
)

// ParkingServer for grpc
type ParkingServer struct {
	parkingService *service.ParkingService
	protocol.UnimplementedParkingServiceServer
}

// NewParkingServer parking lot server
func NewParkingServer(parkingService *service.ParkingService) *ParkingServer {
	return &ParkingServer{parkingService: parkingService}
}

// AddParkingLot adding info about all parking lot
func (srv *ParkingServer) AddParkingLot(ctx context.Context, in *protocol.CreateParkingLotRequest) (*protocol.CreateParkingLotResponse, error) {
	parkingLot := model.ParkingLot{Num: int(in.Num), InParking: in.InParking, Remark: in.Remark}
	err := srv.parkingService.Add(ctx, &parkingLot)
	if err != nil {
		log.Warnf("grpc add parking lot: %v", err)
		return nil, fmt.Errorf("grpc add parking lot %w", err)
	}
	res := &protocol.CreateParkingLotResponse{Num: int32(parkingLot.Num), InParking: parkingLot.InParking, Remark: parkingLot.Remark}
	return res, nil
}

// GetAllParkingLot getting info about all parking lots
func (srv *ParkingServer) GetAllParkingLot(ctx context.Context, in *protocol.Empty) (*protocol.GetAllParkingLotsResponse, error) {
	parkingLots, err := srv.parkingService.GetAll(ctx)
	if err != nil {
		log.Warnf("grpc get all parking lots: %v", err)
		return nil, fmt.Errorf("grpc get parking lots %w", err)
	}
	var resParkingLots []*protocol.ParkingLot

	for i := 0; i < len(parkingLots); i++ {
		resParkingLots[i] = &protocol.ParkingLot{Num: int32(parkingLots[i].Num), InParking: parkingLots[i].InParking, Remark: parkingLots[i].Remark}
	}
	res := &protocol.GetAllParkingLotsResponse{ParkingLot: resParkingLots}
	return res, nil
}

// GetParkingLot getting info about parking lot
func (srv *ParkingServer) GetParkingLot(ctx context.Context, in *protocol.GetParkingLotRequest) (*protocol.GetParkingLotResponse, error) {
	parkingLot, err := srv.parkingService.GetByNum(ctx, int(in.Num))
	if err != nil {
		log.Warnf("grpc get parking lot: %v", err)
		return nil, fmt.Errorf("grpc get parking lot %w", err)
	}
	res := &protocol.GetParkingLotResponse{Num: int32(parkingLot.Num), InParking: parkingLot.InParking, Remark: parkingLot.Remark}
	return res, nil
}

// UpdateParkingLot updating info about parking lot
func (srv *ParkingServer) UpdateParkingLot(ctx context.Context, in *protocol.UpdateParkingLotRequest) (*protocol.UpdateParkingLotResponse, error) {
	err := srv.parkingService.Update(ctx, int(in.Num), in.InParking, in.Remark)
	if err != nil {
		log.Warnf("grpc update parking lot: %v", err)
		return nil, fmt.Errorf("grpc update parking lot %w", err)
	}
	res := &protocol.UpdateParkingLotResponse{Num: in.Num, InParking: in.InParking, Remark: in.Remark}
	return res, nil
}

// DeleteParkingLot deleting info about parking lot
func (srv *ParkingServer) DeleteParkingLot(ctx context.Context, in *protocol.DeleteParkingLotRequest) (*protocol.Empty, error) {
	err := srv.parkingService.Delete(ctx, int(in.Num))
	if err != nil {
		log.Warnf("grpc delete parking lot: %v", err)
		return nil, fmt.Errorf("grpc delete parking lot %w", err)
	}
	return &protocol.Empty{}, nil
}
