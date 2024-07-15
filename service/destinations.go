package service

import (
	pb "content-service/generated/destination"
	"content-service/generated/user"
	"content-service/storage/postgres"
	"context"
	"log/slog"
)

type DestinationService struct {
	DestinationRepo *postgres.DestinationRepo
	UserClient      *user.AuthServiceClient
	Logger          *slog.Logger
}

func (s *DestinationService) ListTravelDestnations(ctx context.Context, in *pb.ListDetinationRequest) (*pb.ListDetinationResponse, error) {
	resp, err := s.DestinationRepo.GetDestinations(in)
	if err != nil {
		s.Logger.Error("sayohat manzillarni ro'yxatini olish")
		return nil, err
	}
	return resp, nil
}

func (s *DestinationService) GetTravelDestination(ctx context.Context, in *pb.GetDestinationRequest) (*pb.GetDestinationResponse, error) {
	resp, err := s.DestinationRepo.GetTravelDestination(in.Id)
	if err != nil {
		s.Logger.Error("Sayohat manzili haqaida malumot olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	return resp, nil
}
